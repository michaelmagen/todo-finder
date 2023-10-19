package todoFinder

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/denormal/go-gitignore"
	"github.com/spf13/viper"
)

// TODO: Handle multi-line comments
type Todo struct {
	text        string
	fileName    string
	relFilePath string
	absFilePath string
	lineNumber  int
}

func (t Todo) Text() string        { return t.text }
func (t Todo) FileName() string    { return t.fileName }
func (t Todo) LineNumber() string  { return strconv.Itoa(t.lineNumber) }
func (t Todo) FilePath() string    { return t.relFilePath }
func (t Todo) FilterValue() string { return t.relFilePath }

func GetTodos(dir string) ([]Todo, error) {
	files, err := getFilesFromDir(dir, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve files from directory: %v", err)
	}

	todos := make([]Todo, 0)
	for _, f := range files {
		todosInFile, err := findTodoInFile(dir, f)
		if err != nil {
			log.Printf("Error while processing file %s: %v", f, err)
			continue // Skip this file and continue with the next one
		}
		todos = append(todos, todosInFile...)
	}
	return todos, nil
}

func findTodoInFile(wd, path string) ([]Todo, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %v", path, err)
	}
	defer file.Close()

	//       TODO:      Explore alternative to buffer scan that can take more characters in file.
	// This works but is inefficient
	// FROM HERE SHOULD BE NEW SMALLER FUNCTION THAT CAN BE CALLED WITH A DIFFERENT CAPACITIES FOR BUFFER
	scanner := bufio.NewScanner(file)
	//const maxCapacity = 10000000
	//buf := make([]byte, maxCapacity)
	//scanner.Buffer(buf, maxCapacity)

	todos := make([]Todo, 0)
	lineNumber := 1

	// This string is what indicates the line is a todo
	todoMarker := viper.GetString("marker")

	for scanner.Scan() {
		// If the line is a comment, add to the todo slice
		lineText := scanner.Text()
		pattern := fmt.Sprintf(`^\s*(\/\/|#|--)\s*%s`, regexp.QuoteMeta(todoMarker))
		re := regexp.MustCompile(pattern)
		// TODO: Export this part in to a seperate function. No more than 3 indents !!!
		if re.MatchString(lineText) {
			lineText = strings.SplitN(lineText, todoMarker, 2)[1]
			lineText = strings.TrimSpace(lineText)
			fileName := filepath.Base(file.Name())
			absPath := file.Name()
			relPath, err := filepath.Rel(wd, absPath)
			if err != nil {
				log.Printf("Failed to get relative path for file %s: %v", absPath, err)
				continue // Skip this line and continue with the next one
			}

			newTodo := Todo{lineText, fileName, relPath, absPath, lineNumber}
			todos = append(todos, newTodo)
		}
		lineNumber++
	}

	if err := scanner.Err(); err != nil {
		// TODO: Modify this so that if it fails, we call the function again with a new buffer with larger size. Maybe make another function that passes in the open file to try to do this.
		return nil, fmt.Errorf("error while scanning file %s: %v", path, err)
	}

	return todos, nil
}

func getFilesFromDir(path string, ignoreMatchers []gitignore.GitIgnore) ([]string, error) {
	dirEntrys, err := os.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory %s: %v", path, err)
	}
	// Get flags from viper
	noGitignore := viper.GetBool("no-gitignore")
	shouldSearchHidden := viper.GetBool("hidden")

	if !noGitignore {
		// Get gitignore file matchers for this current dir
		gitignoreMatcher, err := getGitignoreMatcher(path)

		// If found .gitignore in current dir, add it to the matchers array
		if err == nil {
			ignoreMatchers = append(ignoreMatchers, gitignoreMatcher)
		}
	}

	// Also make sure that files are taking whole path, not just base name
	var files []string

OuterLoop:
	for _, entry := range dirEntrys {
		if entry.IsDir() {
			// Do not explore special directory that start with . like .git
			if strings.HasPrefix(entry.Name(), ".") && !shouldSearchHidden {
				continue
			}

			pathToSubDir := filepath.Join(path, entry.Name())
			// Check that file is not in gitignore
			for _, ignoreMatcher := range ignoreMatchers {
				if ignoreMatcher.Ignore(pathToSubDir) {
					continue OuterLoop
				}
			}

			newFiles, err := getFilesFromDir(pathToSubDir, ignoreMatchers)
			if err != nil {
				log.Printf("Error while exploring directory %s: %v", pathToSubDir, err)
				continue // Skip this directory and continue with the next one
			}
			files = append(files, newFiles...)
		} else {
			pathToFile := filepath.Join(path, entry.Name())
			// Check that file is not in gitignore
			for _, ignoreMatcher := range ignoreMatchers {
				if ignoreMatcher.Ignore(pathToFile) {
					continue OuterLoop
				}
			}
			files = append(files, pathToFile)
		}
	}
	return files, nil
}

func getGitignoreMatcher(path string) (gitignore.GitIgnore, error) {
	// Find .gitignore file
	gitignoreFile := filepath.Join(path, ".gitignore")

	// match a file against a particular .gitignore
	ignore, err := gitignore.NewFromFile(gitignoreFile)
	if err != nil {
		return nil, err
	}

	return ignore, nil
}

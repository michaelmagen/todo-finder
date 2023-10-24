package todoFinder

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/denormal/go-gitignore"
	"github.com/spf13/viper"
)

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

	todos := make([]Todo, 0)
	lineNumber := 1
	// Get todo comment marker from viper
	todoMarker := viper.GetString("marker")

	reader := bufio.NewReader(file)
	// Keep reading lines in until get an End of File error
	for {
		lineText, err := reader.ReadString('\n')
		if err != nil && err != io.EOF {
			return nil, fmt.Errorf("error while reading file %s: %v", path, err)
		}

		if isTodoComment(lineText, todoMarker) {
			// Remove todo marker and whitespace from todo text
			lineText = strings.SplitN(lineText, todoMarker, 2)[1]
			lineText = strings.TrimSpace(lineText)
			fileName := filepath.Base(file.Name())
			absPath := file.Name()
			relPath, err := filepath.Rel(wd, absPath)
			if err != nil {
				log.Printf("Failed to get the relative path for file %s: %v", absPath, err)
				continue // Skip this line and continue with the next one
			}

			newTodo := Todo{lineText, fileName, relPath, absPath, lineNumber}
			todos = append(todos, newTodo)
		}

		// Break from loop since reached end of file
		if err == io.EOF {
			break
		}

		lineNumber++
	}

	return todos, nil
}

func isTodoComment(lineText, todoMarker string) bool {
	// Check if the line is a comment with the todo marker
	pattern := fmt.Sprintf(`^\s*(\/\/|#|--)\s*%s`, regexp.QuoteMeta(todoMarker))
	re := regexp.MustCompile(pattern)
	return re.MatchString(lineText)
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

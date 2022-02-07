package file

import (
	"bufio"
	"errors"
	"io"
	"lazy-log/utils"
	"os"

	"github.com/spf13/cobra"
)

type AnalyzeCommand struct {
	Filepath      string
	SearchPattern string
}

func BuildAnalyzeCommand(cmd *cobra.Command, args []string) (AnalyzeCommand, error) {
	if len(args) != 1 {
		return AnalyzeCommand{}, errors.New("filepath is required as an argument")
	}

	search := cmd.Flag("search").Value.String()

	if search == "" {
		return AnalyzeCommand{}, errors.New("Search text is empty")
	}

	return AnalyzeCommand{
		Filepath:      args[0],
		SearchPattern: search,
	}, nil
}

func CheckIfValidFile(filepath string) (bool, error) {
	if _, err := os.Stat(filepath); err != nil && os.IsNotExist(err) {
		return false, errors.New("Log file doesn't exist")
	}

	return true, nil
}

func LoadWholeFileToMemory(filepath string) string {
	file, err := os.Open(filepath)

	utils.Check(err)

	defer file.Close()

	fileInfo, err := file.Stat()

	utils.Check(err)

	// Get the file size upfront
	fileSize := fileInfo.Size()

	// Init a buffer large enough to hold the content of this file
	buffer := make([]byte, fileSize)

	_, err = file.Read(buffer)

	utils.Check(err)

	return string(buffer)
}

// Load files based on bufferSize
// Usage:
//  LoadFileInChunks(analyzeCommand.Filepath, 100, func(buffer string) {
//   fmt.Println(">>>>>>>> Chunk: ", buffer)
// })
func LoadFileInChunks(filepath string, bufferSize int, processor func(buffer string)) {
	file, err := os.Open(filepath)

	utils.Check(err)

	defer file.Close()

	buffer := make([]byte, bufferSize)

	for {
		bytesRead, err := file.Read(buffer)

		if err != nil && err != io.EOF {
			utils.ExitGracefully(err)
		}

		if err == io.EOF {
			break
		}

		processor(string(buffer[:bytesRead]))
	}
}

func ScanFile(filepath string, processor func(line string)) {
	file, err := os.Open(filepath)

	utils.Check(err)

	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		processor(scanner.Text())
	}
}

func ProcessLogFile(analyzeCommand AnalyzeCommand, processor func(line string)) {
	ScanFile(analyzeCommand.Filepath, processor)
}

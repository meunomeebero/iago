package pkg

import (
	"flag"
	"fmt"
	"os"

	"github.com/robertokbr/iago/utils"
)

func printHelp() {
	fmt.Println("Iago is a CLI for OpenAI's GPT API")
	fmt.Println("Usage: iago [flags]")
	fmt.Println("Flags:")
	fmt.Println("  -help\t\tShow this help message")
	fmt.Println("  -pk\t\tSet your OpenAI API Secret Key")
	fmt.Println("  -prompt\tPrompt to send to OpenAI")
}

func writePK(pk string) {
	path := fmt.Sprintf("%s/.env.iago", utils.GetPathToEnv())

	f, err := os.Create(path)

	if err != nil {
		fmt.Println(err)
		return
	}

	defer f.Close()

	_, err = f.WriteString("OPEN_API_SK=" + pk)

	if err != nil {
		fmt.Println(err)
		f.Close()
		return
	}

	err = f.Close()

	if err != nil {

		fmt.Println(err)
		return
	}

	fmt.Println("OpenAI API Secret Key set")
}

func printSetPK() {
	fmt.Println("Please set your OpenAI API Secret Key")
	fmt.Println("You can find it at https://beta.openai.com/account/api-keys")
	fmt.Println("You can also set it with the -pk flag")
	return
}

func ExecuteCLI() {
	help := flag.Bool("help", false, "Show help")
	pk := flag.String("pk", "", "OpenAI API Secret Key")
	prompt := flag.String("prompt", "", "Prompt to send to OpenAI")
	ytb := flag.Bool("ytb", false, "use youtube data")
	comments := flag.Bool("cmt", false, "use comments data")
	subs := flag.Bool("sub", false, "use subscribers data")

	flag.Parse()

	if *help {
		printHelp()
		return
	}

	if *pk != "" {
		writePK(*pk)
		return
	}

	if utils.GetSK() == "" {
		printSetPK()
		return
	}

	if *prompt != "" {
		text := *prompt
		var data interface{}

		if *ytb {
			if *subs {
				data = GetSubscribersComments()
			}

			if *comments {
				data = GetChannelComments()
			}

			res := AnswerQuestion(
				fmt.Sprintf("%v", data),
				text,
			)

			fmt.Println(res)
			return
		}

		res := AnswerQuestion(text)
		fmt.Println(res)
		return
	}
}

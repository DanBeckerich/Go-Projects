package main

import (
	"bufio"
	"context"
	"encoding/csv"
	"fmt"
	"github.com/BlueMonday/go-scryfall"
	"os"
	"strings"
	"time"
)

func main() {

	wd, err := os.Getwd()
	args := os.Args[1]
	infile, err := os.Open( wd + "\\" + args)
	defer infile.Close()
	if err != nil{
		panic(err)
	}

	outfile, err := os.Create( args + ".csv")
	defer outfile.Close()
	if err != nil{
		panic(err)
	}

	scanner := bufio.NewScanner(infile)
	csvwriter := csv.NewWriter(outfile)

	ctx := context.Background()
	client, err := scryfall.NewClient()
	if err != nil {
		panic(err)
	}

	sco := scryfall.SearchCardsOptions{
		Unique:        scryfall.UniqueModePrints,
		Order:         scryfall.OrderSet,
		Dir:           scryfall.DirDesc,
		IncludeExtras: false,
	}

	err = csvwriter.Write([]string{"QTY", "Card Name", "Card Color", "Rarity", "Set Preference", "Other Set OK?", "STORE USE ONLY"})
	if err != nil {
		panic(err)
	}

	for scanner.Scan() {

		text := strings.ReplaceAll(scanner.Text(), "\r\n", "")

		if text != "" {
			tempSlice := strings.Split(scanner.Text(), " ")

			cardName := strings.Join(tempSlice[1:]," ")

			fmt.Println("fetching", cardName)

			temp, err := client.SearchCards(ctx,cardName , sco)
			if err != nil {
				panic(err)
			}

			Color := "CL"

			if strings.Contains(temp.Cards[0].TypeLine, "Artifact") {
				Color = "A"
			} else if strings.Contains(temp.Cards[0].TypeLine, "Land") {
				Color = "L"
			} else if temp.Cards[0].Colors == nil {
				Color = "C"
			} else if len(temp.Cards[0].Colors) > 1 {
				Color = "Gold"
			} else if len(temp.Cards[0].Colors) == 1{
				Color = string(temp.Cards[0].Colors[0])
 			} else {
 				Color = "ERROR"
			}


				err = csvwriter.Write([]string{"1", temp.Cards[0].Name, Color , temp.Cards[0].Rarity, temp.Cards[0].Set, "Y", ""})
			if err != nil {
				panic(err)
			}

			time.Sleep(200)
		}
	}

	csvwriter.Flush()
	fmt.Println("Finished")
}

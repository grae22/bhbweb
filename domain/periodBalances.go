package domain

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type PeriodBalances struct {
	PeriodId           string
	BalanceByAccountId map[int]AccountBalances
}

type AccountBalances struct {
	Opening Money
	Closing Money
}

func newPeriodBalances(periodId string) PeriodBalances {
	return PeriodBalances{
		PeriodId:           periodId,
		BalanceByAccountId: make(map[int]AccountBalances),
	}
}

func loadPeriodBalances(
	bookName string,
	periodId string,
	accounts accountService,
) (pbals PeriodBalances, subsequentPbals []PeriodBalances, err error) {

	log.Printf("Loading period balances for book \"%s\" and period \"%s\"...", bookName, periodId)

	balancesFilename := bookName + "_" + periodId + ".balances"

	pbals, err = loadPeriodBalancesFromFile(balancesFilename)

	if errors.Is(err, os.ErrNotExist) {
		log.Println("No period balances file, creating one.")

		files, err := getPeriodBalancesFilesOldestFirst(bookName)
		if err != nil {
			return pbals, subsequentPbals, err
		}

		lastFilename := ""

		if len(files) > 0 {
			lastFilename = files[len(files)-1]
			lastPeriodId := lastFilename[strings.Index(lastFilename, "_")+1 : strings.Index(lastFilename, ".")]

			if periodId < lastPeriodId {
				lastFilename = ""
			}
		}

		if lastFilename != "" {
			pbals, err = loadPeriodBalancesFromFile(lastFilename)

			if err != nil {
				return pbals, subsequentPbals, err
			}

			pbals.PeriodId = periodId

			for accId, bals := range pbals.BalanceByAccountId {
				acc, ok := accounts.accountForId(accId)
				if !ok {
					log.Println("Account not found.")
					continue
				}

				accountType := acc.AccountType

				if accountType == TypeBank ||
					accountType == TypeDebtor ||
					accountType == TypeCreditor {

					bals.Opening = bals.Closing
				} else {
					bals.Opening = 0
					bals.Closing = 0
				}

				pbals.BalanceByAccountId[accId] = bals
			}
		} else {
			pbals = newPeriodBalances(periodId)
		}
	} else if err != nil {
		log.Println("Error while reading file.")
		return pbals, subsequentPbals, err
	}

	files, err := getPeriodBalancesFilesOldestFirst(bookName)
	if err != nil {
		return pbals, subsequentPbals, err
	}

	var foundPeriodBalance bool

	for _, f := range files {
		pb, err := loadPeriodBalancesFromFile(f)
		if err != nil {
			return pbals, subsequentPbals, err
		}

		if f == balancesFilename {
			foundPeriodBalance = true
			continue
		}

		if foundPeriodBalance {
			subsequentPbals = append(subsequentPbals, pb)
		}
	}

	return pbals, subsequentPbals, nil
}

func (pb PeriodBalances) save(bookName string) error {
	bytes, err := json.MarshalIndent(pb, "", "  ")
	if err != nil {
		log.Printf("Error serialising period balances (%s).\n", pb.PeriodId)
		return err
	}

	err = os.WriteFile(bookName+"_"+pb.PeriodId+".balances", bytes, os.ModePerm)

	return err
}

func loadPeriodBalancesFromFile(filename string) (pbals PeriodBalances, err error) {
	log.Printf("Loading period balances file \"%s\"...", filename)

	bytes, err := os.ReadFile(filename)

	if err == os.ErrNotExist {
		log.Println("No file found.")
		return
	} else if err != nil {
		log.Println("Error while reading file.")
		return
	}

	log.Println("Deserialising file...")
	err = json.Unmarshal(bytes, &pbals)

	return
}

func getPeriodBalancesFilesOldestFirst(bookName string) ([]string, error) {
	files := []string{}

	err := filepath.Walk(
		".",
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				log.Println("Error while finding period balances files.")
				return err
			}

			if strings.HasPrefix(info.Name(), bookName) &&
				strings.HasSuffix(info.Name(), ".balances") {

				files = append(files, info.Name())
			}

			return nil
		})

	if err != nil {
		return files, err
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i] < files[j]
	})

	return files, nil
}

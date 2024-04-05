package main

import (
	"fmt"
	gt7 "github.com/snipem/go-gt7-telemetry/lib"
	"strings"
	"time"
)

func flush() error {
	_, err := fmt.Print("\033[H\033[2J")
	return err
}

func generateNumberBanner(number int) string {
	digits := []string{
		`
   ##  
  #  # 
  #  # 
  #  # 
  #  # 
   ##  
       `,
		`
   #   
  ##   
   #   
   #   
   #   
  ###  
       `,
		`
  ###  
 #   # 
    #  
   #   
  #    
 ####  
       `,
		`
  ###  
 #   # 
     # 
   ##  
     # 
 #   # 
  ###  
       `,
		`
    #  
   ##  
  # #  
 #  #  
 #####
    #  
    #  
       `,
		`
  ###  
 #     
 #     
  ###  
     # 
 #   # 
  ###  
       `,
		`
  ###  
 #     
 #     
  ###  
 #   # 
 #   # 
  ###  
       `,
		`
  #### 
     # 
    #  
   #   
  #    
  #    
  #    
       `,
		`
  ###  
 #   # 
 #   # 
  ###  
 #   # 
 #   # 
  ###  
       `,
		`
  ###  
 #   # 
 #   # 
  #### 
     # 
 #   # 
  ###  
       `,
	}

	// Convert the number to a string
	numStr := fmt.Sprintf("%d", number)

	// Generate the banner for each digit in the number
	var result []string
	for _, digit := range numStr {
		digitVal := int(digit - '0')
		if digitVal >= 0 && digitVal <= 9 {
			result = append(result, digits[digitVal])
		}
	}

	// Join the digits together

	// Combine the banners for each digit with line breaks
	return strings.Join(mergeLines(result), "\n")

}

func mergeLines(array []string) []string {
	var merged []string

	for _, element := range array {
		lines := strings.Split(element, "\n")
		merged = append(merged, lines...)
	}

	return merged
}

func main() {
	gt7c := gt7.NewGT7Communication("255.255.255.255")
	go gt7c.Run()
	for true {
		_ = flush()
		//fmt.Println("\r" + banner.Inline(fmt.Sprintf("%.0f", gt7c.LastData.CarSpeed)))
		fmt.Println("\r" + generateNumberBanner(int(gt7c.LastData.CarSpeed)))
		fmt.Println("\r" + (fmt.Sprintf("%.0f", gt7c.LastData.CarSpeed)))
		time.Sleep(16 * time.Millisecond)
	}
}

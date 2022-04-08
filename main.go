package main

import (
	"bufio"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
)

type Ascii struct {
	Display string
	Color   string
	Size    string
}

func Index(s string, pattern string) int { //Index is a function that look for pattern inside s

	lenghtS := len(s)
	lenghtPattern := len(pattern)
	for i := 0; i <= lenghtS-lenghtPattern; i++ {
		if s[i:i+lenghtPattern] == pattern {
			return i
		}
	}
	return -1
}

//Transfert_file_To_String is a function that transfer a data from a file into an array of string Line By Line
func Transfert_File_To_String(font string) []string {
	ch := font + ".txt"
	// we open the file so we can manipulate it
	file, err := os.Open(ch)
	//we verify if the file exists
	if err != nil {
		println("Error")

	}

	if err != nil {
		log.Fatalf("")
	}
	//we create a new scanner
	fileScanner := bufio.NewScanner(file)
	//we use Split for the scan
	fileScanner.Split(bufio.ScanLines)
	var fileTextLines []string // this the string array which will contain all the date from the file line by line
	//we scan for next token and we append to fileTextLines
	for fileScanner.Scan() {
		fileTextLines = append(fileTextLines, fileScanner.Text())
	}

	file.Close() //we close the file
	return fileTextLines
}

//this function will Display the string sent in argument in a graphical style
func Ascii_art(ch string, tch []string) string {
	s := ""
	//j := 0
	//each letter has a height of 8 so we display all letters line by line
	for line := 0; line < 8; line++ {
		i := 0 //i index for ch
		for range ch {
			a := ((int(ch[i]) - 32) * 9) + 1 // a represents the position of the letter in the file

			//s[j] = s[j] + (tch[a+line]) //we print the first line then 2nd etc...
			s = s + (tch[a+line]) //we print the first line then 2nd etc...
			i++
		}
		//j++
		s = s + "\n"
	}
	return s
}

//this function is useful to calculate how many "\n" we have in a string
func Calc_retour(s string) int {
	i := 0  // index for s
	nb := 0 //number of "\n" in s
	for range s {
		if s[i] == '\n' {
			nb++
		}
		i++
	}
	return nb
}

//Verif's role is to return true if the character exists in the file standard.txt or its a "\n"
func Verif(s string) bool {
	i := 0
	for (i < len(s)) && ((s[i] >= 32 && s[i] <= 126) || (s[i] == 10 || s[i] == 13)) { //we continue on verifying until we reach the end or we find an invalid character
		i++
	}
	return i >= len(s)
}

//Affiche display the characters in graphical style but it manages the Errors and prints "\n" if needed etc...
func Affiche(font, ch string) string {
	if ch == "" {
		return "the result will be displayed here !!!"
	}
	tch := Transfert_File_To_String(font) // we form an array of string
	chf := "\n"
	if (Verif(ch)) == false { //we verify if its a good string with valid characters or we print "Error"
		print("Error\n")
		return "Error"
	} else {
		if Index(ch, "\n") != -1 { // if we have a "\n" in  a string we have to manage it and go to the next line
			nb := Calc_retour(ch) // number of "\n" in ch
			i := 0                // a counter
			for i <= nb {
				if Index(ch, "\n") != -1 {
					s := ch[:Index(ch, "\n")] //we take the slice before the "\n" and we display it in a graphical style
					chf = chf + Ascii_art(s, tch)
					ch = ch[Index(ch, "\n")+1:] // we take the rest of the string
				} else {
					chf = chf + Ascii_art(ch[0:], tch) //if we dont have any "\n" left then we display the rest of ch
				}
				chf = chf + "\n"
				i++
			}
		} else {
			chf = chf + Ascii_art(ch, tch) // if there isnt any problem and there isnt a "\n" we display ch in a graphical style
		}

	}
	return chf
}
func main() {
	var font string
	var color string
	var text string
	var output string
	var size string
	var t []string
	tmpl, _ := template.ParseGlob("./template/*.html")
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		//Call to ParseForm makes form fields available.
		err := r.ParseForm()
		if err != nil {
			// Handle error here via logging and then return
		}
		size = r.PostFormValue("size")
		output = r.PostFormValue("output")
		color = r.PostFormValue("color")
		font = r.PostFormValue("font")
		text = r.PostFormValue("texte")
		t = strings.Split(text, "\r")
		tch := strings.Join(t, "  ")
		data := Ascii{
			Display: Affiche(font, tch),
			Color:   color,
			Size:    size,
		}

		if output != "" && data.Display != "Error" {
			ch := output + ".txt"
			file, err := os.Create(ch)
			defer file.Close()
			if err != nil {
				print("Error\n")
				return
			}
			file.WriteString(data.Display)
		}
		tmpl.ExecuteTemplate(w, "index", data)
	})
	fs := http.FileServer(http.Dir("assets/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.ListenAndServe(":8080", nil)
}

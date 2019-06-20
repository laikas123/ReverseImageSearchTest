//this is algorithm attempt #1
//it works by taking squares of some size of the image
//we will say .01% size squares to start 
//then it gets the most prominent pixel in that square 

package main

import (
    "fmt"
    "golang.org/x/net/html"
    "github.com/laikas123/downloader"
    "image"
    "image/png"
    "image/jpeg"
    "io"
    "log"
    "math"
    "net/http"
    "net/url"
    "os"
    "path/filepath"
    "strings"
    "sync"
    
)

var (
    	
    tokensLinks = make(chan struct{}, 500)
    tokensImages = make(chan struct{}, 10)
    
    mu sync.Mutex
)

type pixelItem struct {

	red int
	green int
	blue int


}

var myBool bool = true

var closer chan bool

var keepCrawling bool

var fileChannel chan string 

var deleteChannel chan bool




func main() {
	
	 

	fileChannel := make(chan string)
	deleteChannel := make(chan bool)
	fileToCrawlWith := os.Args[1]

	firstLink := []string{"https://en.wikipedia.org/wiki/Bluebird"}
		
	array1 := pixelarray(fileToCrawlWith)
	array3 := pixelarray(fileToCrawlWith)
	filedir, _ := os.Open("/home/logan/Desktop/GoLearning/ReverseImageSearch/EmptyFolder")
    	filedir.Chdir()

	closer = make(chan bool)

	imagesCompared := 0

	keepCrawling = true

	go startImageCrawling(firstLink, fileChannel)

	for keepCrawling {
		
		if imagesCompared == 10 {
			close(fileChannel)
			close(deleteChannel)
			close(closer)
			fmt.Println("done searching, program terminating...")
			os.Exit(1)
		}

		//this array gets edited to compare
		//array1 := pixelarray(fileToCrawlWith)
		fileAdded := <-fileChannel

		go deleteFile(fileAdded, deleteChannel)
		
		fmt.Println(fileAdded)

		fmt.Println("METADATA")
		downloader.Meta(fileAdded)
		array2 := pixelarray(fileAdded)

		if array2 == nil {
			
			continue

		}

		//array3 := pixelarray(fileToCrawlWith)

		copy(array3, array1)

		p1 := &array1

		p2 := &array2

		
		p3 := &array3
		
		

		go shiftAndCompare1(p1, p2, 1, 0, p3, 0, 0, 0, closer)
			
		go shiftAndCompare2(p1, p2, 1, 0, p3, 0, 0, 0, closer)
		
		go shiftAndCompare3(p1, p2, 1, 0, p3, 0, 0, 0, closer)

		go shiftAndCompare4(p1, p2, 1, 0, p3, 0, 0, 0, closer)
		
		bool1 := <- closer
		bool2 := <- closer
		bool3 := <- closer
		bool4 := <- closer
		
		
		if bool1 || bool2 || bool3 || bool4 {
			deleteChannel <- false
		}else{
			deleteChannel <- true
		}

		fmt.Println("RECEIVED ONE")
		
		imagesCompared++
		
		//close(closer)
	
	
	
	
	}
	

	//TODO if not a match return the highest closest possible percent that was achieved through 
	//our attempt at image search

}


//This function returns the slice of slices (rows) of pixels, each as a struct containing its RGB values
func pixelarray(filename string) [][]pixelItem {
	   
	    defer func() {
		if p := recover(); p != nil {
			fmt.Println("error with this image checking next")
		}
	    }()
	    
	    image.RegisterFormat("png", "png", png.Decode, png.DecodeConfig)
	    image.RegisterFormat("jpeg", "jpg", jpeg.Decode, jpeg.DecodeConfig)
	    
	    filename = "./" + filename
	    file, err := os.Open(filename)

	    if err != nil {
		fmt.Println("Error: File could not be opened" + filename)
		os.Exit(1)
	    }

	    defer file.Close()
	
	    pixels, err := getPixels(file)

	    

	    if err != nil {
		fmt.Println("Error: Image could not be decoded" + filename)
		//os.Exit(1)
		return nil
	    }

	    return pixels
	 	
	



}


//this function
//returns to the previous function what we stated it returns the array of pixel items
//which are no more than a struct of RGB values

func getPixels(file io.Reader) ([][]pixelItem, error) {

	    img, _, err := image.Decode(file)

	    if err != nil {
		return nil, err
	    }
	    
	    bounds := img.Bounds()

	    width, height := bounds.Max.X, bounds.Max.Y

	    
	    
	    var pixels [][]pixelItem
	    
	    for y := 0; y < height; y++ {
		
		var row []pixelItem
		
		for x := 0; x < width; x++ {
		    r, g, b := returnPixelRGB(img.At(x, y).RGBA())
		    row = append(row, pixelItem{r, g, b})
		    
		}
		pixels = append(pixels, row)
	    }
		
	   // fmt.Printf("%s", len(pixels))
	    //fmt.Printf("%s", len(pixels[0]))

	   // return pixels, nil
		
           

	    //height = 270
	    if len(pixels) > 100 || len(pixels) < 100 {
		
		p := &pixels
		pixels = resizeHeight(p)
		//fmt.Println("OKHERE")
		//fmt.Println(len(pixels))
	    }
		

	    //width = 340 
	    if len(pixels[0]) > 100 || len(pixels[0]) < 100 {
			
		p := &pixels
		pixels = resizeWidth(p)
		//fmt.Println(len(pixels[0]))
		//fmt.Println("OKHERE2")
		//fmt.Println(pixels)

	    } 

	    return pixels, nil
}


func returnPixelRGB(r uint32, g uint32, b uint32, a uint32) (int, int, int) {
	
	    r1 := int(r / 257)
	    g1 := int(g / 257)
	    b1 := int(b / 257)

	    //discard the alpha channel we aren't concerned with it 
	    _ = a 
	   
	    
	    return r1, g1, b1
}









func shiftAndCompare1(pixelArray1, pixelArray2 *[][]pixelItem, cyclesRecursion int, horizontalTracker int,
		resetHorizontal *[][]pixelItem, highestPercent float64, totalPixelsPossible int, initialTotalPixels int,
		send chan <- bool)  {
	
	
	
	pixelArray1C := *pixelArray1
	pixelArray2C := *pixelArray2
	
	
        

	

	totalMatchingPixels := 0 

	maxNumberHorizontalShifts := len(pixelArray1C[0]) - 1
	
	if totalPixelsPossible == 0 {	
		takeAway := getCountOfValueInSlice(pixelArray1C, pixelItem{255, 255, 255})
		//fmt.Println(takeAway)
		totalPixelsPossible = len(pixelArray1C) * len(pixelArray1C[0])
		initialTotalPixels = totalPixelsPossible
		//fmt.Println(totalPixelsPossible)
		totalPixelsPossible = totalPixelsPossible - takeAway	
	}
	

	
	for i := 0; i < len(pixelArray1C); i++ {
			
			cmpr1 := pixelArray1C[i]
			cmpr2 := pixelArray2C[i]
			for j := 0; j < len(cmpr1); j++ {
				if determineIfTwoPixelsAreTheSame(cmpr1[j], cmpr2[j]) {
					totalMatchingPixels++
				}
			}

		}
	var percent float64 = float64(totalMatchingPixels) / float64(totalPixelsPossible)
	
	if percent > highestPercent {
		highestPercent = percent
	}
	
	//fmt.Println(percent)
	if percent > .15 {
		
		fmt.Println("CASE1: The Images Are a Match % =")
		fmt.Println(percent)
		send <- true
		return 

	}
		
	pixelArray1 = &pixelArray1C
	arrayAfterH := shiftHorizontally(pixelArray1, len(pixelArray1C))
	pixelArray1 = &arrayAfterH
	//totalPixelsPossible = totalPixelsPossible - len(pixelArray1)
	horizontalTracker++
	
	
	

	if horizontalTracker == maxNumberHorizontalShifts || horizontalTracker > maxNumberHorizontalShifts {
		
		
		pixelArray1 = resetHorizontal
		
		
				
		arrayAfterV := shiftVertically(pixelArray1, len(pixelArray1C))
		pixelArray1 = &arrayAfterV
		resetHorizontal = &arrayAfterV
		

		//totalPixelsPossible = totalPixelsPossible + (len(pixelArray1) * len(pixelArray1[0])) - len(pixelArray1[0])
		horizontalTracker = 0
	}
	
	cyclesRecursion = cyclesRecursion + 1
	//myString := fmt.Sprintf("%s", cyclesRecursion, totalPixelsPossible)	
	//fmt.Println(myString)
	if cyclesRecursion == initialTotalPixels {
		fmt.Println("CASE1: the images are not a match closest % =")
		fmt.Println(highestPercent)
		send <- false
		//fmt.Println(cyclesRecursion)
		return 

	}

	
	go shiftAndCompare1(pixelArray1, pixelArray2, cyclesRecursion, horizontalTracker, resetHorizontal, highestPercent, totalPixelsPossible, initialTotalPixels, send)
	
	
	return 
	
}

func shiftAndCompare2(pixelArray1, pixelArray2 *[][]pixelItem, cyclesRecursion int, horizontalTracker int,
		resetHorizontal *[][]pixelItem, highestPercent float64, totalPixelsPossible int, initialTotalPixels int,
		send chan <- bool)  {
	
	
	
	pixelArray1C := *pixelArray1
	pixelArray2C := *pixelArray2

        

	

	totalMatchingPixels := 0 

	maxNumberHorizontalShifts := len(pixelArray1C[0]) - 1
	
	if totalPixelsPossible == 0 {	
		takeAway := getCountOfValueInSlice(pixelArray1C, pixelItem{255, 255, 255})
		//fmt.Println(takeAway)
		totalPixelsPossible = len(pixelArray1C) * len(pixelArray1C[0])
		initialTotalPixels = totalPixelsPossible
		//fmt.Println(totalPixelsPossible)
		totalPixelsPossible = totalPixelsPossible - takeAway	
	}
	

	
	for i := 0; i < len(pixelArray1C); i++ {
			
			cmpr1 := pixelArray1C[i]
			cmpr2 := pixelArray2C[i]
			for j := 0; j < len(cmpr1); j++ {
				if determineIfTwoPixelsAreTheSame(cmpr1[j], cmpr2[j]) {
					totalMatchingPixels++
				}
			}

		}
	var percent float64 = float64(totalMatchingPixels) / float64(totalPixelsPossible)
	
	if percent > highestPercent {
		highestPercent = percent
	}
	
	//fmt.Println(percent)
	if percent > .15 {
		
		fmt.Println("CASE2: The Images Are a Match % =")
		fmt.Println(percent)
		send <- true
		return 

	}
		
	pixelArray1 = &pixelArray1C
	arrayAfterH := shiftHorizontallyOpposite(pixelArray1, len(pixelArray1C))
	pixelArray1 = &arrayAfterH
	//totalPixelsPossible = totalPixelsPossible - len(pixelArray1)
	horizontalTracker++
	
	
	

	if horizontalTracker == maxNumberHorizontalShifts || horizontalTracker > maxNumberHorizontalShifts {
		
		
		pixelArray1 = resetHorizontal
		
		
				
		arrayAfterV := shiftVertically(pixelArray1, len(pixelArray1C))
		pixelArray1 = &arrayAfterV
		resetHorizontal = &arrayAfterV
		

		//totalPixelsPossible = totalPixelsPossible + (len(pixelArray1) * len(pixelArray1[0])) - len(pixelArray1[0])
		horizontalTracker = 0
	}
	
	cyclesRecursion = cyclesRecursion + 1
	//myString := fmt.Sprintf("%s", cyclesRecursion, totalPixelsPossible)	
	//fmt.Println(myString)
	if cyclesRecursion == initialTotalPixels {
		fmt.Println("CASE2: the images are not a match closest % =")
		fmt.Println(highestPercent)
		send <- false
		//fmt.Println(cyclesRecursion)
		return 

	}

	
	go shiftAndCompare2(pixelArray1, pixelArray2, cyclesRecursion, horizontalTracker, resetHorizontal, highestPercent, totalPixelsPossible, initialTotalPixels, send)

	return 
	




}


func shiftAndCompare3(pixelArray1, pixelArray2 *[][]pixelItem, cyclesRecursion int, horizontalTracker int,
		resetHorizontal *[][]pixelItem, highestPercent float64, totalPixelsPossible int, initialTotalPixels int,
		send chan <- bool)  {
	
	
	
	pixelArray1C := *pixelArray1
	pixelArray2C := *pixelArray2

        

	

	totalMatchingPixels := 0 

	maxNumberHorizontalShifts := len(pixelArray1C[0]) - 1
	
	if totalPixelsPossible == 0 {	
		takeAway := getCountOfValueInSlice(pixelArray1C, pixelItem{255, 255, 255})
		//fmt.Println(takeAway)
		totalPixelsPossible = len(pixelArray1C) * len(pixelArray1C[0])
		initialTotalPixels = totalPixelsPossible
		//fmt.Println(totalPixelsPossible)
		totalPixelsPossible = totalPixelsPossible - takeAway	
	}
	

	
	for i := 0; i < len(pixelArray1C); i++ {
			
			cmpr1 := pixelArray1C[i]
			cmpr2 := pixelArray2C[i]
			for j := 0; j < len(cmpr1); j++ {
				if determineIfTwoPixelsAreTheSame(cmpr1[j], cmpr2[j]) {
					totalMatchingPixels++
				}
			}

		}
	var percent float64 = float64(totalMatchingPixels) / float64(totalPixelsPossible)
	
	if percent > highestPercent {
		highestPercent = percent
	}
	
	//fmt.Println(percent)
	if percent > .15 {
		
		fmt.Println("CASE3: The Images Are a Match % =")
		fmt.Println(percent)
		send <- true
		return 

	}
		
	pixelArray1 = &pixelArray1C
	arrayAfterH := shiftHorizontally(pixelArray1, len(pixelArray1C))
	pixelArray1 = &arrayAfterH
	//totalPixelsPossible = totalPixelsPossible - len(pixelArray1)
	horizontalTracker++
	
	
	

	if horizontalTracker == maxNumberHorizontalShifts || horizontalTracker > maxNumberHorizontalShifts {
		
		
		pixelArray1 = resetHorizontal
		
		
				
		arrayAfterV := shiftVerticallyOpposite(pixelArray1, len(pixelArray1C))
		pixelArray1 = &arrayAfterV
		resetHorizontal = &arrayAfterV
		

		//totalPixelsPossible = totalPixelsPossible + (len(pixelArray1) * len(pixelArray1[0])) - len(pixelArray1[0])
		horizontalTracker = 0
	}
	
	cyclesRecursion = cyclesRecursion + 1
	//myString := fmt.Sprintf("%s", cyclesRecursion, totalPixelsPossible)	
	//fmt.Println(myString)
	if cyclesRecursion == initialTotalPixels {
		fmt.Println("CASE3: the images are not a match closest % =")
		fmt.Println(highestPercent)
		send <- false
		//fmt.Println(cyclesRecursion)
		return 

	}

	
	go shiftAndCompare3(pixelArray1, pixelArray2, cyclesRecursion, horizontalTracker, resetHorizontal, highestPercent, totalPixelsPossible, initialTotalPixels, send)

	return 
	

}


func shiftAndCompare4(pixelArray1, pixelArray2 *[][]pixelItem, cyclesRecursion int, horizontalTracker int,
		resetHorizontal *[][]pixelItem, highestPercent float64, totalPixelsPossible int, initialTotalPixels int,
		send chan <- bool)  {
	
	
	
	pixelArray1C := *pixelArray1
	pixelArray2C := *pixelArray2

        

	

	totalMatchingPixels := 0 

	maxNumberHorizontalShifts := len(pixelArray1C[0]) - 1
	
	if totalPixelsPossible == 0 {	
		takeAway := getCountOfValueInSlice(pixelArray1C, pixelItem{255, 255, 255})
		//fmt.Println(takeAway)
		totalPixelsPossible = len(pixelArray1C) * len(pixelArray1C[0])
		initialTotalPixels = totalPixelsPossible
		//fmt.Println(totalPixelsPossible)
		totalPixelsPossible = totalPixelsPossible - takeAway	
	}
	

	
	for i := 0; i < len(pixelArray1C); i++ {
			
			cmpr1 := pixelArray1C[i]
			cmpr2 := pixelArray2C[i]
			for j := 0; j < len(cmpr1); j++ {
				if determineIfTwoPixelsAreTheSame(cmpr1[j], cmpr2[j]) {
					totalMatchingPixels++
				}
			}

		}
	var percent float64 = float64(totalMatchingPixels) / float64(totalPixelsPossible)
	
	if percent > highestPercent {
		highestPercent = percent
	}
	
	//fmt.Println(percent)
	if percent > .15 {
		
		fmt.Println("CASE4: The Images Are a Match % =")
		fmt.Println(percent)
		send <- true
		return 

	}
		
	pixelArray1 = &pixelArray1C
	arrayAfterH := shiftHorizontallyOpposite(pixelArray1, len(pixelArray1C))
	pixelArray1 = &arrayAfterH
	//totalPixelsPossible = totalPixelsPossible - len(pixelArray1)
	horizontalTracker++
	
	
	

	if horizontalTracker == maxNumberHorizontalShifts || horizontalTracker > maxNumberHorizontalShifts {
		
		
		pixelArray1 = resetHorizontal
		
		
				
		arrayAfterV := shiftVerticallyOpposite(pixelArray1, len(pixelArray1C))
		pixelArray1 = &arrayAfterV
		resetHorizontal = &arrayAfterV
		

		//totalPixelsPossible = totalPixelsPossible + (len(pixelArray1) * len(pixelArray1[0])) - len(pixelArray1[0])
		horizontalTracker = 0
	}
	
	cyclesRecursion = cyclesRecursion + 1
	//myString := fmt.Sprintf("%s", cyclesRecursion, totalPixelsPossible)	
	//fmt.Println(myString)
	if cyclesRecursion == initialTotalPixels {
		fmt.Println("CASE4: the images are not a match closest % =")
		fmt.Println(highestPercent)
		send <- false
		//fmt.Println(cyclesRecursion)
		return 

	}

	
	go shiftAndCompare4(pixelArray1, pixelArray2, cyclesRecursion, horizontalTracker, resetHorizontal, highestPercent, totalPixelsPossible, initialTotalPixels, send)

	return 
	




}





//FIXED VERSION

func shiftHorizontally(pixelArray *[][]pixelItem, length int) [][]pixelItem {
	
	nullPixel  := pixelItem{-1, -1, -1}
	
	useDontEdit := *pixelArray
	
	lengthOfInnerSlices := len(useDontEdit[0])

	bufferSlices := make([][]pixelItem, length, length)

	for o := 0; o < length; o++ {
		
		bufferSlices[o] = make([]pixelItem, lengthOfInnerSlices, lengthOfInnerSlices)

	}

	for i := 0; i < length; i++ {
		
		bufferArrayInner := bufferSlices[i]

		useDontEditInner := useDontEdit[i]
		
		for z := 0; z < lengthOfInnerSlices; z++ {
			if z == 0 {
				appendThis := nullPixel
				bufferArrayInner[z] = appendThis
			}else{
			
			appendThis := useDontEditInner[z - 1]
			
			bufferArrayInner[z] = appendThis
			
			}
		}

		bufferArrayInner = bufferArrayInner[:len(bufferArrayInner) - 1]

	}


	return bufferSlices
	
}


func shiftHorizontallyOpposite(pixelArray *[][]pixelItem, length int) [][]pixelItem {
	
	nullPixel  := pixelItem{-1, -1, -1}
	

	
	useDontEdit := *pixelArray
	

	
	lengthOfInnerSlices := len(useDontEdit[0])

	 
	bufferSlices := make([][]pixelItem, length, length)


	for o := 0; o < length; o++ {
		
		bufferSlices[o] = make([]pixelItem, lengthOfInnerSlices, lengthOfInnerSlices)

	}



	for i := 0; i < length; i++ {
		
		bufferArrayInner := bufferSlices[i]
		
		useDontEditInner := useDontEdit[i]

		for z := 0; z < lengthOfInnerSlices; z++ {
			if z == lengthOfInnerSlices - 1 {
							
				appendThis := nullPixel
				bufferArrayInner[z] = appendThis
			}else{
			
			appendThis := useDontEditInner[z + 1]
			
			bufferArrayInner[z] = appendThis
			
			}
		}
		
		bufferArrayInner = bufferArrayInner[:len(bufferArrayInner) - 1]

	}

	return bufferSlices
	
}




func shiftVertically(pixelArray *[][]pixelItem, length int) [][]pixelItem {
	
	nullPixel  := pixelItem{-1, -1, -1}

	useDontEdit := *pixelArray

	lengthOfInnerSlices := len(useDontEdit[0])

	bufferSlices := make([][]pixelItem, length + 1, length + 1)

	for o := 0; o < length; o++ {
	
		bufferSlices[o] = make([]pixelItem, lengthOfInnerSlices, lengthOfInnerSlices)
	}

	for i := 0; i < length; i++ {
		
		bufferArrayInner := bufferSlices[i]

		useDontEditInner := useDontEdit[i]

		for z := 0; z < len(useDontEditInner); z++ {
			
			appendThis := useDontEditInner[z]

			bufferArrayInner[z] = appendThis

		}

	}

	nullSliceRow := make([]pixelItem, lengthOfInnerSlices, lengthOfInnerSlices)
	for u := 0; u < lengthOfInnerSlices; u++ {
		nullSliceRow[u] = nullPixel
	}
	
	bufferSlices[len(bufferSlices) - 1] = nullSliceRow
	


	bufferSlices = bufferSlices[1:]


	return bufferSlices
	
}


func shiftVerticallyOpposite(pixelArray *[][]pixelItem, length int) [][]pixelItem {
	
	nullPixel  := pixelItem{-1, -1, -1}

	useDontEdit := *pixelArray

	lengthOfInnerSlices := len(useDontEdit[0])



	bufferSlices := make([][]pixelItem, length + 1, length + 1)
	
	//bufferslices is one row longer than the original for appending the -1 value

	for o := 0; o < length; o++ {
		
		bufferSlices[o] = make([]pixelItem, lengthOfInnerSlices, lengthOfInnerSlices)

	}

	
	nullSliceRow := make([]pixelItem, lengthOfInnerSlices, lengthOfInnerSlices)
	for u := 0; u < lengthOfInnerSlices; u++ {
		nullSliceRow[u] = nullPixel
	}
	
	bufferSlices[0] = nullSliceRow
	for i := 1; i < length; i++ {
		
		bufferArrayInner := bufferSlices[i]
		
		useDontEditInner := useDontEdit[i - 1]
		
		for z := 0; z < len(useDontEditInner); z++ {
			
			appendThis := useDontEditInner[z]

			bufferArrayInner[z] = appendThis

		}

	}
	

	bufferSlices = bufferSlices[:len(bufferSlices) - 1]


	
	

	

	return bufferSlices
	
}







func getCountOfValueInSlice(pixelArray [][]pixelItem, valueOfInterest pixelItem) int {
	instancesOfValOfInterest := 0
	for i := 0; i < len(pixelArray); i++ {
	
		cmprArray := pixelArray[i]
		for j := 0; j < len(cmprArray); j++ {
			if cmprArray[j] == valueOfInterest {
				instancesOfValOfInterest++
			}
		}
	}

	return instancesOfValOfInterest
	
}



func determineIfTwoPixelsAreTheSame(pixelItem1, pixelItem2 pixelItem) bool {


	r1 := pixelItem1.red
	g1 := pixelItem1.green
	b1 := pixelItem1.blue

	r2 := pixelItem2.red
	g2 := pixelItem2.green
	b2 := pixelItem2.blue

	match := true

	if ((math.Abs(float64(r1) - float64(r2))) < float64(30)) && r1 != -1 && r2 != -1 && r1 != 255 && r2 != 255 {
		
	}else{
		match = false
	}
	if (math.Abs(float64(g1) - float64(g2))) < float64(30) && g1 != -1 && g2 != -1 && g1 != 255 && g2 != 255  {
		
	}else{
		match = false
	}
	if (math.Abs(float64(b1) - float64(b2))) < float64(30) && b1 != -1 && b2 != -1 && b1 != 255 && b2 != 255 {
		
	}else{
		match = false
	}

	return match
	
}















const arrayLength = 3

func resizeHeight(pixelArray *[][]pixelItem) [][]pixelItem {




	myArray := *pixelArray
	
	p := &myArray

	finalResult := resize(p)
	
	

	return finalResult

}



func resize(pixelArray *[][]pixelItem)  [][]pixelItem{


	
	compareDontEdit := *pixelArray

	
	
	//270
	bufferSlices := make([][]pixelItem, len(compareDontEdit), len(compareDontEdit))
	
		
	//give the array empty rows 
	for o := 0; o < len(bufferSlices); o++ {
		
		bufferSlices[o] = make([]pixelItem, len(compareDontEdit[0]), len(compareDontEdit[0]))

	}

	//for 270 copy all the original data into our new array
	for i := 0; i < len(compareDontEdit); i++ {
		bufferArrayInner := bufferSlices[i]

		compareDontEditInner := compareDontEdit[i]
		
		for z := 0; z < len(compareDontEditInner); z++ {
			
			appendThis := compareDontEditInner[z]
			
			bufferArrayInner[z] = appendThis

		}

	}
	//fmt.Println("buffer slices")
	
	//270
	initialHeight := len(bufferSlices)
	
	//270
	height := len(bufferSlices)

	//check for remainder 
	x := 0
	x = height % 100
	remainder := 0
	//fmt.Println(height)



	//there is potentially a special case....
	

	scaled1 := make([][]pixelItem, 0, 0)

	//if it is not divisible by 100
	if x != 0  || height < 100  {
		for {
					
			height = height * 2
			remainder = height % 100
			//fmt.Println(remainder)
			
			if remainder <= 80 && height > 100 {
				//fmt.Println(height)				
				break
			}
			
		}	



		//we end up with height 540 remainder 40 
		
		//scaled1 has len and cap = 580
		scaled1 = make([][]pixelItem, (height + remainder), (height + remainder))
		
		

		//fill scaled one with all null pixels
		for o := 0; o < len(scaled1); o++ {
			
			scaled1[o] = make([]pixelItem, len(bufferSlices[0]), len(bufferSlices[0]))
			inner1 := scaled1[o]		
			for index, _ := range scaled1[o] {
				inner1[index] = pixelItem{-1, -1, -1}
			}

		}		

		

		loopsThrough := 0

		
		for i := 0; i < len(bufferSlices); i++ {

			rowToScale := bufferSlices[i]

			//fmt.Println(rowToScale)
			//neither will this
			lengthRow := len(rowToScale)

			//fmt.Println("1loop")
			
			
			for w := 0; w < ((height + remainder) / initialHeight); w++ {
				//get the scaled row of

				
				//fmt.Println("2loop")
		
				innerScaled1 := scaled1[loopsThrough ]

				//fmt.Println("2loop")
				
				for z := 0; z < lengthRow; z++ {
					//fmt.Println("3loop")
					appendThis := rowToScale[z]

										
					//fmt.Println(appendThis)
					
					innerScaled1[z] = appendThis

					//fmt.Println(innerScaled1)
					

				}
				loopsThrough++
			}
			

		}


		//fmt.Println("OK")
		//fmt.Println(len(scaled1))
		
		checkLength := len(scaled1)
		
		p := &scaled1

		indices := returnEveryOtherIndex(p)
		
		

		lengthIndices := len(indices)

		//p2 := &indices

		//reduce := reduceAllByOne(p2)

		//fmt.Println(reduce)
		//fmt.Println(checkLength)
		//fmt.Println("CHECKLENGTH")
		for checkLength > 100 {
			
				p := &scaled1

				indices := returnEveryOtherIndex(p)	

				lengthIndices = len(indices)

				
				
				for i := 0; i < lengthIndices; i++ { 
					pointerArray := &scaled1
				
					scaled1 = trimSlice(pointerArray, indices[i])
					pointerIndices := &indices
					indices = reduceAllByOne(pointerIndices)
					
					if (len(scaled1)) == 100 {
						
						
						pointerArray = &scaled1				
						fixInt := fixArray(pointerArray)

						if(fixInt == -1){
							return scaled1
						}
				
						emptyAtTheEnd := len(scaled1[fixInt:])
						//fmt.Printf("%s", "fixInt", fixInt)
						//fmt.Printf("%s", "adjustFromHere", emptyAtTheEnd)

				
				
						//indentBy := ((10 - fixInt) / 2) + 1

				


						//fmt.Println(appendThis[fixInt:])


						blankAppend := scaled1[fixInt: fixInt + (emptyAtTheEnd/2)]
						//fmt.Println("OK")
						appendThisSpecific := scaled1[:fixInt]

						//blankOnEitherEnd := appendThis[fixInt:indentBy]

						scaled1 = append(blankAppend, appendThisSpecific...)
						scaled1 = append(scaled1, blankAppend...)
						
						

						blankPixelRow := make([]pixelItem, len(scaled1[0]), len(scaled1[0]))

						for index, _ := range blankPixelRow {
							blankPixelRow[index] = pixelItem{-1, -1, -1}
						}

						for len(scaled1) < 100 {
					
							scaled1 = append(scaled1, blankPixelRow)
						
						}
						return scaled1
						
					}
				}

		
				checkLength = len(scaled1)

			}
	}
		
	
	return make([][]pixelItem, 0, 0)
	

}







func trimSlice(pixelArray *[][]pixelItem, index int) [][]pixelItem {

	compareDontEdit := *pixelArray

	

	//given some random index let us say 3 which is value 4 
	//we want to remove that value and get a slice with all the original values in the same order 
	//without the 4 

	//we will let z = the index to remove 

	if index == 54 {
		index = 53
	}
	
	z := index	


	sliceNew1 := compareDontEdit[:z]

	sliceNew2 := compareDontEdit[z+1:]

	sliceNew1 = append(sliceNew1, sliceNew2...)

	compareDontEdit = sliceNew1

	

	return compareDontEdit

}


func returnEveryOtherIndex(pixelArray *[][]pixelItem) []int {

	compareDontEdit := *pixelArray
	
	extraPad := 0

	if (len(compareDontEdit) & 2) == 1 {
		extraPad = 1
	}


	indexes := make([]int, (len(compareDontEdit)/2 + extraPad), (len(compareDontEdit)/2 + extraPad))
	//fmt.Println("INDEXLENGTH")
	//fmt.Println(len(indexes))
	//indexes has length 145
	iterateVal := len(compareDontEdit)
	//fmt.Printf("%s", "myLength", len(compareDontEdit))
	//580	
	if  iterateVal < 200 {
		
		iterateVal = iterateVal - 2

	}
	
	//fmt.Printf("%s", "iterationValue", iterateVal)

	//fmt.Println("\n\n")
	for i := 0; i < iterateVal; i ++ {
		//fmt.Println(i)
		
		if ((i == 0) || (i % 2 == 0)) && (i/2) < len(indexes)   {
			
			indexes[i/2] = i
		}	


	}
	//fmt.Println("wrongPlace")
	return indexes
	

}

func reduceAllByOne(intSlice *[]int) []int {

	compareDontEdit := *intSlice

	
	
	for i := 0; i < len(compareDontEdit); i ++ {

		val := compareDontEdit[i]

		val = val - 1

		compareDontEdit[i] = val	


	}

	return compareDontEdit
	

}



func fixArray(pixelArray *[][]pixelItem)  int{


	compareDontEdit := *pixelArray

	for index, item := range compareDontEdit {
		
		if item[0].red == -1 {
			return index
		}

	}

	return -1 

}









func resizeWidth (pixelArray *[][]pixelItem) [][]pixelItem{
	
	myArray := *pixelArray
	
	p := &myArray

	finalResult := resizeW(p)

	//fmt.Println("FINALRESULT")

	//fmt.Println(len(finalResult[0]))

	//fmt.Println(len(finalResult))

	return finalResult

}


func resizeW(pixelArray *[][]pixelItem) [][]pixelItem {
	
	compareDontEdit := *pixelArray
	

	//width original is going to be 7
	width := len(compareDontEdit[0])
	widthOriginal := len(compareDontEdit[0])
	remainder := 0
		

	if width < 100 || width > 100 {
		
		for {
			width = width*2
			remainder = width % 100
			
			if width > 100 && remainder < 80 {
				//fmt.Println(width)
				//fmt.Println(remainder)				
				break
			}

		    }


	}

	copyInSize := make([][]pixelItem, len(compareDontEdit), len(compareDontEdit))

	for o := 0; o < len(copyInSize); o++ {
		

		copyInSize[o] = make([]pixelItem, (width + remainder), (width + remainder))
		innerCopy := copyInSize[o]
		for q := 0; q < len(innerCopy); q++ {
			innerCopy[q] = pixelItem{-1, -1, -1}
		}
	}	


	

	//14 / 7 = 2 
	numberToScaleBy := (width + remainder) / widthOriginal 

	//fmt.Printf("%s", numberToScaleBy)
	
	loopsThrough := 0
	
	//fmt.Println(len(compareDontEdit))


	//this runs 4 times 
	for l := 0; l < len(compareDontEdit); l++ {



		
		innerCompareDontEdit := compareDontEdit[l]
		innerCopyInSize := copyInSize[l]
		
		 
		loopsThrough = 0
		for u := 0; u < len(innerCompareDontEdit) ; u++ {
				appendThis := innerCompareDontEdit[u]
				//fmt.Println("OK2")
				
				//this runs twice 
				//a needs to keep adding ie 0, 1, 2, 3, 4, 5
				for a := 0; a < numberToScaleBy; a++ {
					//fmt.Println(a + loopsThrough)
					innerCopyInSize[a + loopsThrough] = appendThis	
					//fmt.Println("OK3")		

			}

			loopsThrough = loopsThrough + 2
			
		}
			
		//fmt.Println(copyInSize)


	}
	
	//note every time through the loop 
	//variable assignments are new 
	//to use a variable throughout a loop without it changing it has to be assigned outside of
	//the loop of interest

	//continueToNext := false

	//breakCompletely := false
	

	var appendThis []pixelItem




	for len(copyInSize[0]) > 100 {

	for p := 0; p < len(copyInSize) ; p++ {
		//fmt.Println("OKSTILL")
		appendThis = copyInSize[p]
		
		pointerArray := &appendThis
		
		indices := returnEveryOtherIndexW(pointerArray)
		pointerIndices := &indices

		lengthIndices := len(indices)
			
		for t := 0; t < lengthIndices; t++ {
			pointerArray = &appendThis
			appendThis  = trimRowW(pointerArray, indices[t])
			pointerIndices = &indices
			indices = reduceAllByOneW(pointerIndices)
			//fmt.Println("FINEHERE")
			if len(appendThis) == 100 {
				//fmt.Println("FINEHERESTILL")
				//fmt.Println(appendThis)
				pointerArray = &appendThis				
				fixInt := fixRowW(pointerArray)

				if(fixInt == -1) {
					break
				}
				
				emptyAtTheEnd := len(appendThis[fixInt:])
				//fmt.Printf("%s", "FixInt", fixInt)
				//fmt.Printf("%s", "EmptyAtTheEnd", emptyAtTheEnd)
				
		
				//fmt.Println(appendThis[fixInt:])


				blankAppend := appendThis[fixInt: fixInt + (emptyAtTheEnd/2)]
				//fmt.Println("OK")
				appendThisSpecific := appendThis[:fixInt]



				appendThis = append(blankAppend, appendThisSpecific...)
				appendThis = append(appendThis, blankAppend...)
				
				for len(appendThis) < 100 {
					//fmt.Println(appendThis
					appendThis = append(appendThis, pixelItem{-1, -1, -1})
				
				}
				//fmt.Printf("%s", "length", len(appendThis))
				break
			}
		} 

		copyInSize[p] = appendThis

	
		}

	}

	return copyInSize

}



func trimRowW(pixelArray *[]pixelItem, index int) []pixelItem {

	compareDontEdit := *pixelArray

	

	//given some random index let us say 3 which is value 4 
	//we want to remove that value and get a slice with all the original values in the same order 
	//without the 4 

	//we will let z = the index to remove 
	
	z := index	

	sliceNew1 := compareDontEdit[:z]

	sliceNew2 := compareDontEdit[z+1:]

	sliceNew1 = append(sliceNew1, sliceNew2...)

	compareDontEdit = sliceNew1

	//fmt.Println("allgoodHere")

	return compareDontEdit

}


func returnEveryOtherIndexW(pixelArray *[]pixelItem) []int {

	compareDontEdit := *pixelArray

	extraPad := 0 
	
	if (len(compareDontEdit) % 2 == 0) {
		//extraPad = 1
	}

	indexes := make([]int, (len(compareDontEdit)/2) + extraPad, (len(compareDontEdit)/2) + extraPad)
	
	for i := 0; i < len(compareDontEdit); i ++ {

		if ((i == 0) || (i % 2 == 0)) && (i/2) < len(indexes) {
			indexes[i/2] = i
		}	


	}

	return indexes
	

}

func reduceAllByOneW(intSlice *[]int) []int {

	compareDontEdit := *intSlice

	
	
	for i := 0; i < len(compareDontEdit); i ++ {

		val := compareDontEdit[i]

		val = val - 1

		compareDontEdit[i] = val	


	}

	return compareDontEdit
	

}




func fixRowW(pixelRow *[]pixelItem)  int{


	compareDontEdit := *pixelRow 

	for index, item := range compareDontEdit {
		
		if item.red == -1 {
			return index
		}

	}

	return -1 

}




func startImageCrawling(sliceStrings []string, filechannel chan <- string) {

	
	//make sense of it from here 
	
	//the worklist is where new links get sent to 	

	//in hindsight setting this to 1 doesn't help all that much 

	//the way we are iterating through this is in a for loop iterating through a slice which is constantly

	//getting added to so one slice could mean 10000 links 

	//perhaps the worklist should be of type string

	worklist := make(chan []string)
	
	//we start out work list in a new goroutine by sending our initial 
	//list of links, aka the slice defined in main passed here
	go func() { worklist <- sliceStrings }()

	// a map that prevents duplicate links from being crawled again
	seen := make(map[string]bool)
	
	//a channel to 
	//
	


	//we call range on a channel which is essentially equivalent to an attempt to receive from a channel
	//we sent one request already so the first receive executes fine 
	//in our case this was only one in length
	continueIsOk := true
	//The takeaway, this essentially will try to receive forever until the channel is closed 
	
	


	//attempts to receive until channel is close 
	for list := range worklist {
		if len(tokensLinks) >= 500 && len(tokensImages) >= 10{
				continueIsOk = false
				fmt.Println("break1")
				break
		}
		//iterates through each string presented by the worklist's slices of strings 
		for _, link := range list {
			if len(tokensLinks) >= 500 && len(tokensImages) >= 1{
				continueIsOk = false
				fmt.Println("break2")
				break
			}

			//only crawl links that have not been seen and mark them as seen once we do 
			if !seen[link] && continueIsOk {
				seen[link] = true
				//crawl the link and return the result to the worklist				
				go func(link string) {
					
					//the reason this works is because
										
					worklist <- crawl(link, filechannel)
					
					
		
				}(link)
				
			}
		}
	}

	//draining the tokens and worklist channel to prevent leak
	for  {
		_, open := <- worklist
		if !open {
			close(worklist)
			break
		}
		
	}
	for  {
		_, open := <- tokensLinks
		if !open {
			close(tokensLinks)
			break
		}
		
	}
	for  {
		_, open := <- tokensImages
		if !open {
			close(tokensImages)
			break
		}
		
	}

	
	fmt.Println("channels closed OK")

	//close(tokensLinks)
	//close(tokensImages)
	//close(worklist)


}



func crawl(url string, filechannel chan <- string) []string {
	fmt.Println(url)


	//crawl returns a slice of strings which represent links added to the worklist 
	
	//crawl calls extract for every url for crawl when it is called 

	//lets limit the number of calls to extract 

	
	list, err := extract(url, filechannel)
	
	if err != nil {
		log.Print(err)
	}

	

	return list
}



func extract(url string, filechannel chan <- string) ([]string, error) {
	

			
	//connect to the url
	resp, err := http.Get(url)
	

	//if err occurs fine 
	if err != nil {
		return nil, err
	}
	//if error occurs fine
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("getting %s: %s", url, resp.Status)
	}
	//this parses the body portion of an html document
	doc, err := html.Parse(resp.Body)
	//save the body to a variable then close it 	
	resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("parsing %s as HTML: %v", url, err)
	}
	

	//this is the slice we return 
	var links []string

	i := 0
	l := 0 

	breakLinks := false

	breakImages := false

	//gets a node from the html document
	visitNode := func(n *html.Node) {

		if len(tokensImages) >= 10 {
			
			breakImages = true
		
		}

		if len(tokensLinks) >= 500 {
			
			breakLinks = true
		
		}
		
		
		if n.Type == html.ElementNode && n.Data == "img" && len(tokensImages) < 10 {
			for _, img := range n.Attr {
				//fmt.Println("runningImages")
				//href is an attribute of html anchor that means a link is embedded 
				if img.Key == "src" {
					imageString, err := resp.Request.URL.Parse(img.Val)
					extensionType := imageString.String()
					extensionType = extensionType[len(extensionType) - 3:]
					//fmt.Println(extensionType)
					if extensionType != "png" && extensionType != "jpg" {
						//don't deal with SVG images						
						continue 
					}
					
					
					if err != nil {
						continue // ignore bad URLs
					}
					
					if !breakImages {
					
						filenamewhencreated := imagedownloaderupdated(imageString.String())
											
						tokensImages <- struct{}{}
						//fmt.Println("fileMakesItBackHere" + filenamewhencreated)
						filechannel <- filenamewhencreated
					
						fmt.Println(filenamewhencreated)
					
					}
					
				}
				
				
				
				

			}
		}

	
		//only check for links if we don't have 5 new ones 
		if n.Type == html.ElementNode && n.Data == "a" && len(tokensLinks) < 500 {
			for _, newLink := range n.Attr {
				//fmt.Println("runningLinks")								
				if len(links) >= 5 && i >= 10  {
					fmt.Println("broke")					
					break
				}
				if newLink.Key == "href" {
					addThis, err := resp.Request.URL.Parse(newLink.Val)
					
					//ignore bad links
					if err != nil {
						continue
					}
					l++

					if !breakLinks {
						
			
						tokensLinks <- struct{}{}
						links = append(links, addThis.String())
						
					}
				}
				
			}
			
		}	
		
				
		
		
	}

	
	
	
	
	forEachNode(doc, visitNode, nil)
	return links, nil

	
}


//this recursively visits every node from the body of a website 


func forEachNode(n *html.Node, pre, post func(n *html.Node)) {
	if pre != nil {
		//go func() {tokensChan <- struct{}{} }()
		pre(n)
		
	}
		
	//What I want to implement is say once the worklist has 5 links 
	//get all the images from those links 	
	//then proceed with all the remaining receives 
	//and let the image compare do its work 

	//this is the real magic it goes to each node 
	//by recursively calling forEachNode on the nodes 
	//this gaurantees that we search the whole pages nodes 
	for c := n.FirstChild; c != nil; c = c.NextSibling {		
		forEachNode(c, pre, post)
	}
	
	
	if post != nil {
		post(n)
	}
}





func imagedownloaderupdated(fileurl string) string{
	

    //change dir to empty folder
    filedir, _ := os.Open("/home/logan/Desktop/GoLearning/ReverseImageSearch/EmptyFolder")
    filedir.Chdir()

    fullUrlFile := fileurl
    
    // Build fileName from fullPath
    fileName := buildFileName(fullUrlFile)

    // Create blank file
    file := createFile(fileName)

    // Put content on file
    putFile(file, httpClient(), fullUrlFile, fileName)
	
    mu.Lock()
	if len(file.Name()) > 0 {
		
	}else{
	
	}
    mu.Unlock()
    	
	

    return file.Name()

}




func putFile(file *os.File, client *http.Client, fullUrlFile string, fileName string) {
    mu.Lock()
    defer mu.Unlock()
    resp, err := client.Get(fullUrlFile)
    path, _ := filepath.Abs(filepath.Dir(file.Name()))
    fmt.Println(path + "/" + file.Name())

    filedir, _ := os.Open("/home/logan/Desktop/GoLearning/ReverseImageSearch/EmptyFolder")
    filedir.Chdir()
    

    checkError(err)

    defer resp.Body.Close()

    size, err := io.Copy(file, resp.Body)

    defer file.Close()

    checkError(err)

    fmt.Println("Just Downloaded a file %s with size %d", fileName, size)
    
}

func buildFileName(fullUrlFile string) string{
    fileUrl, err := url.Parse(fullUrlFile)
    checkError(err)

    path := fileUrl.Path
   
    segments := strings.Split(path, "/")

    fileName := segments[len(segments)-1]

    return fileName
}

func httpClient() *http.Client {
    client := http.Client{
        CheckRedirect: func(r *http.Request, via []*http.Request) error {
            r.URL.Opaque = r.URL.Path
            return nil
        },
    }

    return &client
}

func createFile(fileName string) *os.File {
    file, err := os.Create(fileName)
   
    checkError(err)
    return file
}

func checkError(err error) {
    defer func() {
		if p := recover(); p != nil {
			fmt.Println("error with this image checking next")
		}
	    }()
    if err != nil {
       panic(err)
    }
}


func deleteFile(fileName string, deleteChan <- chan bool){

	deleteBool := <- deleteChan 
	
	if !deleteBool {
		return 
	}else{
		err := os.Remove(fileName)
		if err != nil {
			fmt.Println("error deleting file, skipping")
			return
		}
		
		
	}
	


}







































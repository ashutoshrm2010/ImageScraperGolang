package service

import (
	"log"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"io"
	"github.com/pborman/uuid"
	"gopkg.in/mgo.v2"
	"github.com/sourcecode/ImageScrapGolang/model"
	"encoding/json"
	"gopkg.in/mgo.v2/bson"
	"image/color"
	"image"
	"fmt"
	"image/jpeg"
	"image/png"

	"math"
	"os"
)
//This function is to search image from google and download into in local with jpg extension
func ImageScrapfromGoogleService(userInput string) ([]byte,error) {

	doc, err := goquery.NewDocument("https://www.google.co.in/search?q="+userInput+"&client=ubuntu&hs=Pnm&channel=fs&dcr=0&source=lnms&tbm=isch&sa=X&ved=0ahUKEwjIt-WQ3aXXAhXIfrwKHS1BCj8Q_AUICigB&biw=1301&bih=671")
	if err != nil {
		log.Fatal(err)
	}
	var saveUrl []string
	doc.Find("img").Each(func(index int, item *goquery.Selection) {
		linkTag := item
		link, _ := linkTag.Attr("src")
		//linkText := linkTag.Text()
		//fmt.Printf("Link #%d: '%s' - '%s'\n", index, linkText, link)
		//allLinks=append(allLinks,link)

		saveUrl=append(saveUrl,DownloadImage(link))
	})
	SaveUrlsInMongo(userInput,saveUrl)
	//ChangeImage(saveUrl[0])
	response := make(map[string]interface{})
	response["message"]="success"
	finalResponse, _ := json.Marshal(response)
	return finalResponse, nil
}

func DownloadImage(url string)string{
	response, e := http.Get(url)
	if e != nil {
		log.Fatal(e)
	}
	defer response.Body.Close()
	var filePath string
	filePath="/home/ashu/Desktop/ashu/"+uuid.New()+".jpg"

	file, err := os.Create(filePath)
	if err != nil {
		log.Fatal(err)
	}
	_, err = io.Copy(file, response.Body)
	if err != nil {
		log.Fatal(err)
	}
	file.Close()
	return filePath
}

//saving all image urls into mongo
func SaveUrlsInMongo (searchKey string,urls []string){
	var mongoInsertInput model.MongoInsert
	mongoInsertInput.ID=uuid.New()
	mongoInsertInput.SearchKey=searchKey
	mongoInsertInput.SaveUrl=urls
	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	defer session.Close()
	c := session.DB("webScraping").C("Images")
	c.Insert(&mongoInsertInput)
}

//listing all the inputs given by users
func ListUserSearchInputsService () ([]byte,error) {
	var mongoInsertInput []model.ListSearchKey
	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	defer session.Close()
	c := session.DB("webScraping").C("Images")
	c.Find(nil).Select(bson.M{"SearchKey":1}).All(&mongoInsertInput)

	response := make(map[string]interface{})
	response["userInputs"]=mongoInsertInput
	finalResponse, _ := json.Marshal(response)
	return finalResponse, nil
}
//listing all the urls from mongo based on input given by user
func GetSearchedImageUrlsFromDBService (id string) ([]byte,error) {
	var mongoInsertInput []model.MongoInsert
	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	defer session.Close()
	c := session.DB("webScraping").C("Images")
	c.Find(bson.M{"_id":id}).All(&mongoInsertInput)

	response := make(map[string]interface{})
	response["imageUrls"]=mongoInsertInput
	finalResponse, _ := json.Marshal(response)
	return finalResponse, nil
}

//This Part is for filter image in another format needs time for completing
//don't consider this code as now
//logics are for saving images with different extension and applying grey scale filter


type Converted struct {
	Img image.Image
	Mod color.Model
}

func (c *Converted) At(x, y int) color.Color{
	return c.Mod.Convert(c.Img.At(x,y))
}
//
func ChangeImage(img1 io.Reader,imageFormat string) string{

	//imageFormat,_:=guessImageFormat(img1)
	fmt.Println("change",imageFormat)
	img, _, err := image.Decode(img1)
	if err != nil {
		log.Fatalln("sdsddsdsd",err)
	}
	bounds := img.Bounds()
	w, h := bounds.Max.X, bounds.Max.Y
	grayScale := image.NewGray(image.Rectangle{image.Point{0, 0}, image.Point{w, h}})
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			imageColor := img.At(x, y)
			rr, gg, bb, _ := imageColor.RGBA()
			r := math.Pow(float64(rr), 2.2)
			g := math.Pow(float64(gg), 2.2)
			b := math.Pow(float64(bb), 2.2)
			m := math.Pow(0.2125*r+0.7154*g+0.0721*b, 1/2.2)
			Y := uint16(m + 0.5)
			grayColor := color.Gray{uint8(Y >> 8)}
			grayScale.Set(x, y, grayColor)
		}
	}
	imagePathJPEG:="/home/ashu/Desktop/ashu/test/"+uuid.New()+".jpeg"
	imagePathPNG:="/home/ashu/Desktop/ashu/test/"+uuid.New()+".png"


	fmt.Println("tttttttttt",imageFormat)
	if imageFormat=="jpeg"{
		fmt.Println("jpeggggggggggg")
		outfile, err := os.Create(imagePathJPEG)
		if err != nil {
			log.Fatalln("aiii",err)
		}
		defer outfile.Close()

		err=jpeg.Encode(outfile,grayScale,nil)
		if err!=nil{
			fmt.Println("aaaaaaaaaaa",err)
		}

		return imagePathJPEG

	}else if imageFormat=="png"{
		fmt.Println("pngggggggggggggg",imageFormat)
		outfile, err := os.Create(imagePathPNG)
		if err != nil {
			log.Fatalln("aiii",err)
		}
		defer outfile.Close()

		err=png.Encode(outfile,grayScale)
		if err!=nil{
			fmt.Println("aaaaaaaaaaa",err)
		}

		return imagePathPNG

	}

	/*img, _, err := image.Decode(img1)
	if err != nil {
		log.Fatalln(err)
	}
	bounds := img.Bounds()
	w, h := bounds.Max.X, bounds.Max.Y
	grayScale := image.NewGray(image.Rectangle{image.Point{0, 0}, image.Point{w, h}})
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			imageColor := img.At(x, y)
			rr, gg, bb, _ := imageColor.RGBA()
			r := math.Pow(float64(rr), 2.2)
			g := math.Pow(float64(gg), 2.2)
			b := math.Pow(float64(bb), 2.2)
			m := math.Pow(0.2125*r+0.7154*g+0.0721*b, 1/2.2)
			Y := uint16(m + 0.5)
			grayColor := color.Gray{uint8(Y >> 8)}
			grayScale.Set(x, y, grayColor)
		}
	}*/
	////  imagepath:="/home/ashu/Desktop/ashu/test/"+uuid.New()+".jpeg"
	// imagepath1:="/home/ashu/Desktop/ashu/test/"+uuid.New()+".png"


	/*if imageFormat=="jpeg"||imageFormat=="jpg"{
		outfile, err := os.Create(imagepath)
		if err != nil {
			log.Fatalln("aiii",err)
		}
		defer outfile.Close()
		err=jpeg.Encode(outfile,grayScale,nil)
		if err!=nil{
			fmt.Println("aaaaaaaaaaa",err)
		}

		return imagepath
       }else if imageFormat=="png" {
		outfile, err := os.Create(imagepath1)
		if err != nil {
			log.Fatalln("aiii",err)
		}
		defer outfile.Close()
		err=png.Encode(outfile,grayScale)
		if err!=nil{
			fmt.Println("aaaaaaaaaaa",err)
		}
		return imagepath1
	}*/
	return ""
}

func guessImageFormat(r io.Reader) (format string, err error) {
	_, format, err = image.DecodeConfig(r)
	return format,err
}
/*
func init() {
	image.RegisterFormat("jpeg", "jpeg", jpeg.Decode, jpeg.DecodeConfig)
	image.RegisterFormat("png", "png", png.Decode, png.DecodeConfig)
	image.RegisterFormat("gif", "gif", gif.Decode, gif.DecodeConfig)

}
*/

/*imageFormat,_:=guessImageFormat(response.Body)
	if imageFormat=="jpeg"{
		filePath="/home/ashu/Desktop/ashu/"+uuid.New()+".jpeg"
		file, err := os.Create(filePath)
		if err != nil {
			log.Fatal(err)
		}
		// Use io.Copy to just dump the response body to the file. This supports huge files
		_, err = io.Copy(file, response.Body)
		if err != nil {k
			log.Fatal(err)
		}
		file.Close()
		return filePath
	}else if imageFormat=="png"{
		filePath="/home/ashu/Desktop/ashu/"+uuid.New()+".png"
		file, err := os.Create(filePath)
		if err != nil {
			log.Fatal(err)
		}
		// Use io.Copy to just dump the response body to the file. This supports huge files
		_, err = io.Copy(file, response.Body)
		if err != nil {
			log.Fatal(err)
		}
		file.Close()
		return filePath
	}*/
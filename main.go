package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"sync"

	"github.com/go-vk-api/vk"
)

var _api = Generator{}

type Generator struct {
	client *vk.Client
}

func (Generator Generator) Init() *vk.Client {
	if Generator.client != nil {
		return Generator.client
	}
	Generator.client, _ = vk.NewClientWithOptions(
		vk.WithToken("TOKEN"),
	)
	return Generator.client
}

type Zaeb struct {
	Count int `json:"count"`
	Items []struct {
		Attachments []struct {
			Photo struct {
				Sizes []struct {
					Url string `json:"url"`
				} `json:"sizes"`
			} `json:"photo"`
		} `json:"attachments"`
	} `json:"items"`
}

func (Generator Generator) GetPhotos(id int) []string {
	var res []string
	var ZaebZaeb Zaeb
	i := 0
	Generator.Init().CallMethod("wall.get", vk.RequestParams{"owner_id": "-" + strconv.Itoa(id)}, &ZaebZaeb)
	for _, item := range ZaebZaeb.Items {
		for _, attach := range item.Attachments {
			if len(attach.Photo.Sizes) > 0 {
				res = append(res, attach.Photo.Sizes[len(attach.Photo.Sizes)-1].Url)
			}
		}
	}
	var deep func(i int)

	var wg sync.WaitGroup
	var rwm sync.RWMutex

	deep = func(i int) {
		var Zaeb Zaeb
		defer wg.Done()
		Generator.Init().CallMethod("wall.get", vk.RequestParams{"owner_id": "-" + strconv.Itoa(id), "offset": i, "count": 100}, &Zaeb)
		for _, item := range Zaeb.Items {
			for _, attach := range item.Attachments {
				if len(attach.Photo.Sizes) > 0 {
					rwm.Lock()
					res = append(res, attach.Photo.Sizes[len(attach.Photo.Sizes)-1].Url)
					rwm.Unlock()
				}
			}
		}
		fmt.Println(strconv.Itoa(i) + "/" + strconv.Itoa(ZaebZaeb.Count))
	}

	for i < ZaebZaeb.Count {
		for j := 0; j < 10; j++ {
			if i < ZaebZaeb.Count {
				wg.Add(1)
				i += 100
				go deep(i)
			}
		}
		wg.Wait()
	}

	return res
}

func (Generator Generator) GetId(str string) int {
	var IdZaeb []struct {
		Id int `json:"id"`
	}
	Generator.Init().CallMethod("groups.getById", vk.RequestParams{"group_id": str}, &IdZaeb)

	me := IdZaeb

	return me[0].Id
}

func main() {
	var id = _api.GetId("digitallow")
	kek := _api.GetPhotos(id)
	html := `<!DOCTYPE html>
	<html lang="en">
	
	<head>
		<meta charset="UTF-8">
		<meta http-equiv="X-UA-Compatible" content="IE=edge">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
	
	</head>
	
	<body style="width: 100vw; margin:0px; height: 100vh; overflow: hidden;">
		<canvas></canvas>

	</body>
	<script>
		document.querySelector("canvas").width = document.body.clientWidth;
		document.querySelector("canvas").height = document.body.clientHeight;
	
		function getRandomInt(max) {
			return Math.floor(Math.random() * Math.floor(max));
		}
	
		let _images = JSON.parse(`
	html += "`"
	json_data2, _ := json.Marshal(kek)
	html += string(json_data2)
	html += "`"

	html += `);
	let images = JSON.parse(localStorage.getItem('images'));
	if(images === null){
		images = JSON.parse(JSON.stringify(_images));
	}
	localStorage.setItem('images', JSON.stringify(images));
	let render = () => {
		images = JSON.parse(localStorage.getItem('images'));
    	let currentId = getRandomInt(images.length - 1);
        if (images.length < 2) {
            images = JSON.parse(JSON.stringify(_images))
        }
        const img = images[currentId];
        images.splice(currentId, 1);
		let canvas = document.querySelector("canvas");
		let ctx = canvas.getContext("2d"); 
		const drawImage = (image, width, height) => {
			ctx.drawImage(image, width, height)
		}
		

		var imgC = new Image();

		imgC.addEventListener('load', () => {
			var width = imgC.naturalWidth; 
			var height = imgC.naturalHeight;
			for (let i = 0; i <= 20; i++) {
                for (let j = 0; j <= 20; j++) {
                    drawImage(imgC, Number(width)*i, Number(height)*j)
				}
			}
			localStorage.setItem('images', JSON.stringify(images));
			setTimeout(() => {
				render()
			}, 2000);
		}, false);

		imgC.src = img;
	}
    render();


</script>

</html>
	`

	ioutil.WriteFile("index.html", []byte(html), 0644)

}

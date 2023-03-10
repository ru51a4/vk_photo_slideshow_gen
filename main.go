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
	var Zaeb Zaeb
	i := 0
	Generator.Init().CallMethod("wall.get", vk.RequestParams{"owner_id": "-" + strconv.Itoa(id)}, &Zaeb)
	for _, item := range Zaeb.Items {
		for _, attach := range item.Attachments {
			if len(attach.Photo.Sizes) > 0 {
				res = append(res, attach.Photo.Sizes[len(attach.Photo.Sizes)-1].Url)
			}
		}
	}
	var deep func(id int)

	var wg sync.WaitGroup

	deep = func(i int) {
		defer wg.Done()
		Generator.Init().CallMethod("wall.get", vk.RequestParams{"owner_id": "-" + strconv.Itoa(id), "offset": i, "count": 100}, &Zaeb)
		for _, item := range Zaeb.Items {
			for _, attach := range item.Attachments {
				if len(attach.Photo.Sizes) > 0 {
					res = append(res, attach.Photo.Sizes[len(attach.Photo.Sizes)-1].Url)
				}
			}
		}
		fmt.Println(strconv.Itoa(i) + "/" + strconv.Itoa(Zaeb.Count))
	}

	for i < Zaeb.Count {
		for j := 0; j < 10; j++ {
			if i < Zaeb.Count {
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
		<title>Document</title>
		<script src="https://pixijs.download/release/pixi.js"></script>
	
	</head>
	
	<body style="width: 100vw; margin:0px; height: 100vh; overflow: hidden;">
	
	</body>
	<script>
		let app = new PIXI.Application({ width: document.body.clientWidth, height: document.body.clientHeight });
		document.body.appendChild(app.view);
	
	
		function getRandomInt(max) {
			return Math.floor(Math.random() * Math.floor(max));
		}
	
		let _images = JSON.parse(`
	html += "`"
	json_data2, _ := json.Marshal(kek)
	html += string(json_data2)
	html += "`"

	html += `);
    let images = JSON.parse(JSON.stringify(_images));
	let counter = 0;
    let render = (width, height) => {
        counter++
		let currentId = getRandomInt(images.length - 1);
        if (images.length < 2) {
            images = JSON.parse(JSON.stringify(_images))
        }
        const img = images[currentId];
        images.splice(currentId, 1);
		let cc = counter
        let sprite = PIXI.Sprite.from(img);
        sprite.texture.baseTexture.on('loaded', () => {
			if(counter !=cc){
				return
			}
            let width = sprite.width
            let height = sprite.height;
            for (let i = 0; i <= 20; i++) {
                for (let j = 0; j <= 20; j++) {
                    let q = PIXI.Sprite.from(img);
                    q.x = i * width
                    q.y = j * height
                    app.stage.addChild(q);
                }
            }
			setTimeout(() => {
				render(width, height)
			}, 2000);
        });
        app.stage.addChild(sprite);
	}
    render();


</script>

</html>
	`

	ioutil.WriteFile("index.html", []byte(html), 0644)

}

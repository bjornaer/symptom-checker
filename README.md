# sympton-checker
![tests](https://github.com/bjornaer/sympton-checker/actions/workflows/push.yaml/badge.svg)[![Go Report Card](https://goreportcard.com/badge/github.com/bjornaer/sympton-checker)](https://goreportcard.com/report/github.com/bjornaer/sympton-checker)

*Cough!* Welcome to Symptom-Checker

This is a service which runs a very rudimentary Gin Gonic server and React UI, whose sole purpose is to provide you with a tool to find ailments that match-up to your health status. **This is by no means a medical diagnosis, if you feel unwell go see a doctor.**

---
### Cloning

do your usual, `git clone blabla` but just to take into account, the `Makefile` assumes you will place the codebase under `~/personal/sympton-checker`
So either place it in the same structure, or go into the makefile and at the begining make sure to replace the `BUILD_DIR` variable to point to the proper directory... *why this* you might ask? Well, cause that is my folder structure, that's why.

---
### Run 

#### Dependencies

To use this piece of software you need to have Docker and docker-compose.

You might want to run it directly on your machine _and that is fine_. Then you need NPM and Go.

Once you have the deps in, you gotta do `go mod download` on the root and `npm install` inside `/fronted` dir.
#### Actually running it
To run this bad boy locally you have three options:

        make local

or

        make docker

or the less sexy but equally cool:

        docker-compose up

navigate to your [localhost](http://localhost:8081) and check it out.

---
### API


The Backend service runs a simple REST api. The endpoints in place are to POST a List of Symptoms, GET a List of Symptoms.

A source data file can be passed in as an ENV var under the name `SYMPTOMS_FILE` which is used to populate our in memory DB.

- `POST /api/symptoms`
- `GET /api/symptoms`

You can run the API by itself by calling:

        make run

#### EXAMPLES
- `POST /api/symptoms`
    ```sh
    curl -X POST --location "http://localhost:8081/api/symptoms" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer whatever" \
    -d  "[\"HP:0000256\", \"HP:0001249\"]"
    ```
- `GET /api/symptoms`
    ```sh
    curl http://localhost:8081/api/symptoms
    ```
---
### UI
The UI consists of a list of symptoms, which you can select by clicking, and a search box from where you can type how you feel and see which results line up with your symptom. The search box performs a fuzzy search so type it corerctyly or incorrectly and it's got your back. Once you finish selecting symptoms click the pretty button in the middle of the screen and get a list of most likely ailments you might have based on your symptoms.

You can click on the ailments to see a link to a related expert to consult and also read on other symptoms related to the ailment.

![Visual aid to description above](https://media.giphy.com/media/y4D5CjpXRhbAi7xjkO/giphy.gif)
---
### Docker

1. run `docker-compose up`
2. ???
3. profit

---
### Tests

To execute tests please refer to the [Makefile](./Makefile) once again and invoke

```sh
make test
```

You can also make a check for _suspicious constructs_ by calling:

```sh
make vet
```

this generates a report file with any findings.

---
### Development

During the development please be sure to run and add tests, as well as formating the code.

If you don't have your favorite code editor already set up for this, you can format the code by running:

```sh
make fmt
```

To clean up any generated file you can run

```sh
make clean
```
---
###  TODO/Author Notes

- I made the app work by each time it spins up it sets a sqlite3 in memory and pulls the data from the internet to populate the DB, of course this would be nice to replace with an actual DB, perhaps psql in this case, and populate it once, have that volume be preserved and be done with it, but for the purposes of the intended usage I left it as is.

- On the UI there are a couple of graphing related components, to make a bar chart/histogram of the diseases and the frequency they came up on the app's search (which is the criteria used to rank them). Anyway, the graph bit is not being rendered because I couldn't get it working properly, but I left the components there for exploring, I would like to enable it eventually.

*Other things to add:*

- Add more testing, I covered the bare minimum -> added a table test for the server
- Missing tests for the front-end
- Setup a CD pipeline when merged to main, or when a release is made.
---

<!-- Improved compatibility of back to top link: See: https://github.com/othneildrew/Best-README-Template/pull/73 -->
<a name="readme-top"></a>
<!--
*** Thanks for checking out the Best-README-Template. If you have a suggestion
*** that would make this better, please fork the repo and create a pull request
*** or simply open an issue with the tag "enhancement".
*** Don't forget to give the project a star!
*** Thanks again! Now go create something AMAZING! :D
-->



<!-- PROJECT SHIELDS -->

<!-- PROJECT LOGO -->
<br />

<h1 align="center">mongoent</h1>



<!-- TABLE OF CONTENTS -->
<details>
  <summary>Table of Contents</summary>
  <ol>
    <li>
      <a href="#about-the-project">About The Project</a>
      <ul>
        <li><a href="#built-with">Built With</a></li>
      </ul>
    </li>
    <li>
      <a href="#getting-started">Getting Started</a>
      <ul>
        <li><a href="#prerequisites">Prerequisites</a></li>
        <li><a href="#installation">Installation</a></li>
      </ul>
    </li>
    <li><a href="#usage">Usage</a></li>
    <li><a href="#roadmap">TodoList</a></li>
    <li><a href="#contributing">Contributing</a></li>
    <li><a href="#license">License</a></li>
  </ol>
</details>



<!-- ABOUT THE PROJECT -->
## About The Project

This project uses CRUD syntax similar to `ent`'s SQL format to invoke the `go.mongodb.org/mongo-driver` package. When the query conditions become too complex, code snippets like the following will be generated in the `mongo-driver`.

```go
filter := bson.D{
		{"$or", bson.A{
			bson.D{
				{"$and", bson.A{
					bson.D{{"age", bson.D{{"$gte", 18}}}},             
					bson.D{{"name", bson.D{{"$regex", "John"}}}},     
				}},
			},
			bson.D{
				{"city", "NY"},                                     
			},
		}},
	}
```

Mongoent utilizes the syntax of the `ent` framework to implement CRUD operations. Here's an example of a query(Please refer to the `cmd/generate_test.go` file for a specific demo.):

```go
newClient := mongoschema.NewClient(mongoschema.Driver(*client))

all, err := newClient.User.SetDBName("my_mongo").Query().
   Where(user.UserNameRegex("c*")).
   All(ctx)
if err != nil {
   log.Fatal(err)
   return
}
```

<p align="right">(<a href="#readme-top">back to top</a>)</p>



### Built With

* golang1.19+

* goimport

* gofumpt

* gofmt

  

<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- GETTING STARTED -->
## Getting Started

### Prerequisites

Make sure you have Go programming language and related tools installed correctly.
### Installation

1.Install gofmt、goimports、gofumpt
   ```sh
   go install golang.org/x/tools/cmd/gofmt@latest
   go install golang.org/x/tools/cmd/goimports@latest
   go install mvdan.cc/gofumpt/gofumports@latest
   ```

<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- USAGE EXAMPLES -->
## Usage

To generate the relevant code, simply run the following command:

```shell
go run -mod=mod github.com/zaihui/mongoent/cmd generate --schemaFile {$YourModelPath} --outputPath {$outputPath} --projectPath {$projectPath} --goModPath {$goModPath}
```

--schemaFile:The schema corresponding to MongoDB is represented as a struct type.

--outputPath：The generated code is typically placed in a directory specified by the project's configuration or convention. Please refer to your project's documentation or configuration files to determine the specific path where the code will be generated.

--projectPath：project path

--goModPath:   go mod path

<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- Todo List -->

- [x] The query conditions typically support the following operators：$eq、$ne、$gt、$lt、$gte、$lte、$regex、$in、$in
- [ ] create operator
- [ ] update operator
- [ ] delete Feature
- [ ] multiple operators

<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- CONTRIBUTING -->
## Contributing

Contributions are what make the open source community such an amazing place to learn, inspire, and create. Any contributions you make are **greatly appreciated**.

If you have a suggestion that would make this better, please fork the repo and create a pull request. You can also simply open an issue with the tag "enhancement".
Don't forget to give the project a star! Thanks again!

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- LICENSE -->
## License

Distributed under the MIT License. See `LICENSE.txt` for more information.

<p align="right">(<a href="#readme-top">back to top</a>)</p>


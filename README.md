# gommand

Gommand provides an easy to use and high performance commands router and processor for Go(lang) using disgord as the base Discord framework. 

## Contributing
Do you have something which you wish to contribute? Feel free to make a pull request. The following simple criteria applies:

- **Is it useful to everyone?:** If this is a domain specific edit, you probably want to keep this as middleware since it will not be accepted into the main project.
- **Does it slow down parsing by >1ms?:** This will likely be denied. We want to keep parsing as high performance as possible.
- **Have you ran `gofmt -w .` and [golangci-lint](https://golangci-lint.run/usage/install/)?:** We like to stick to using Go standards within this project, therefore if this is not done you may be asked to do this for it to be acccepted.

Have you experienced a bug? If so, please make an issue! We take bugs seriously and this will be a large priority for us.

To work on the documentation for gommand, you will want to install [MkDocs](https://www.mkdocs.org/) on your system. From here, you can simply run `mkdocs serve` to serve the documentation on a local web server. When your changes are merged to master, GitHub Actions will automatically deploy your changes to GitHub Pages.

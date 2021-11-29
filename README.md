<div id="top"></div>

<!-- TABLE OF CONTENTS -->
<details>
  <summary>Table of Contents</summary>
  <ol>
    <li>
      <a href="#about-the-project">About Imagswap</a>
    </li>
    <li>
      <a href="#getting-started">Getting Started</a>
      <ul>
        <li><a href="#prerequisites">Prerequisites</a></li>
        <li><a href="#installation">Installation</a></li>
      </ul>
    </li>
    <li><a href="#usage">Usage</a></li>
    <li><a href="#roadmap">Roadmap</a></li>
    <li><a href="#contributing">Contributing</a></li>
    <li><a href="#license">License</a></li>
    <li><a href="#contact">Contact</a></li>
    <li><a href="#acknowledgments">Acknowledgments</a></li>
  </ol>
</details>



<!-- ABOUT THE PROJECT -->
## About The Project
Imageswap is a simple program written in Go that is used to perform image registry hostname transformations as a Kubernetes [mutating webook](https://kubernetes.io/docs/reference/access-authn-authz/admission-controllers/#mutatingadmissionwebhook).
<p align="right">(<a href="#top">back to top</a>)</p>

<!-- GETTING STARTED -->
## Getting Started

This is an example of how you may give instructions on setting up your project locally.
To get a local copy up and running follow these simple example steps.

### Prerequisites
- Go 1.17.1
- Docker
- KinD

### Installation
To compile only the `imageswap` binary:
```
go build -o imageswap main.go
```
Or, to output to `build/_output/bin/`
```
make build
```

Alternatively, to build the binary and create a container image to host it:
```
make build-image
```

<p align="right">(<a href="#top">back to top</a>)</p>



<!-- USAGE EXAMPLES -->
## Usage
Please see `Makefile` and `internal/deploy/` for materials used to test the webhook locally using KinD.  The manifests in `internal/deploy/` can be used as a starting place for deploying

<p align="right">(<a href="#top">back to top</a>)</p>

<!-- ROADMAP -->
## Roadmap
All dates are TBD, this project is passively developed

- [] Makefile test
  - [] Unit Testing
- [] Makefile fmt
- [] Makefile lint
  - [] golangcli-lint
- [] Validating webhook
- [] More complete examples
- [] CI pipeline
- [] Serve on 8080 (insecure)
- [] Helm Chart

See the [open issues](https://github.com/chaospuppy/imageswap/issues) for a full list of proposed features (and known issues).

<p align="right">(<a href="#top">back to top</a>)</p>

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

<p align="right">(<a href="#top">back to top</a>)</p>



<!-- LICENSE -->
## License

Distributed under the Apache 2.0 License. See `LICENSE` for more information.

<p align="right">(<a href="#top">back to top</a>)</p>



<!-- CONTACT -->
## Contact

Tim Seagren - [LinkedIn](https://www.linkedin.com/in/tim-seagren-7876aa112/) - seagrentime@gmail.com

Project Link: [https://github.com/chaospuppy/imageswap](https://github.com/chaospuppy/imageswap)

<p align="right">(<a href="#top">back to top</a>)</p>



<!-- ACKNOWLEDGMENTS -->
## Acknowledgments

* [douglasmakey's admissionscontroller](https://github.com/douglasmakey/admissioncontroller) was used for large portions of the `server`, `pod`, and `hook` packages
* [morvencaos's mutating webhook tutorial](https://github.com/morvencao/kube-mutating-webhook-tutorial) was leveraged for it's Kubernetes CSA and webhook CA patch scripts (`internal/deploy/webhook-create-signed-cert.sh` and `internal/deploy/webhook-patch-ca-bundle.sh` respectively)
* [KinD with registry](https://kind.sigs.k8s.io/docs/user/local-registry/#create-a-cluster-and-registry)

<p align="right">(<a href="#top">back to top</a>)</p>

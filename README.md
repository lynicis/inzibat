<div align="center">
    <p>
        <h1 size="10rem">Inzibat 🪖</h1>
        This word comes from Turkish Language, which means Military Police.
    </p>
</div>

<br/>

<div align="center">
    <strong>Light-weight REST API Gateway</strong>
    <a href="https://github.com/Lynicis/inzibat/actions/workflows/ci.yaml/badge.svg?branch=master&event=push"></a>
    <div>
        <a href="https://github.com/lynicis/inzibat/blob/master/LICENSE">
            <img 
                src="https://img.shields.io/github/license/Lynicis/inzibat"
                alt="License"
            />
        </a>
        <a>
            <img
                src="https://sonarcloud.io/api/project_badges/measure?project=Lynicis_inzibat&metric=coverage"
                alt="SonarCloud Coverage"
            />
        </a>
        <a>
            <img 
                src="https://sonarcloud.io/api/project_badges/measure?project=Lynicis_inzibat&metric=security_rating"
                alt="SonarCloud Security Rating"
            />
        </a>
        <a>
            <img 
                src="https://sonarcloud.io/api/project_badges/measure?project=Lynicis_inzibat&metric=vulnerabilities"
                alt="SonarCloud Vulnerabilities"
            />
        </a>
    </div>
</div>

<br/>

<h2>Installation</h2>

<p>You need install Go 1.20 or higher version, You can download from <a href="https://golang.org/dl/">here</a>.</p>

```bash
go install github.com/Lynicis/inzibat
```

<br/>

<h2>Usage</h2>

```bash
go run inzibat
```

With custom json file:
```bash
go run inzibat -c config.json
```

<br/>

### TODO For v1.0.0:
- [x] Basic Routing
- [ ] Health Check
- [ ] Configurable Logging
- [ ] Configurable Rate Limiting
- [ ] Token-Based Authentication
- [ ] UI

<br/>



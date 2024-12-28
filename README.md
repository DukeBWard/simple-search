```mermaid

graph TB
    %% Styles and classes
    classDef frontend fill:#3498db,stroke:#2980b9,color:white
    classDef backend fill:#2ecc71,stroke:#27ae60,color:white
    classDef search fill:#f1c40f,stroke:#f39c12,color:black
    classDef database fill:#e74c3c,stroke:#c0392b,color:white
    classDef background fill:#9b59b6,stroke:#8e44ad,color:white

    %% Frontend Layer
    subgraph Frontend
        UI["HTMX UI"]:::frontend
        Templates["Templating System (templ)"]:::frontend
        CSS["Tailwind + DaisyUI"]:::frontend
        UI --> Templates
        UI --> CSS
    end

    %% Backend Layer
    subgraph Server
        Fiber["Fiber Web Server"]:::backend
        Routes["Route Handlers"]:::backend
        Auth["JWT Authentication"]:::backend
        Fiber --> Routes
        Fiber --> Auth
    end

    %% Search Engine Core
    subgraph SearchCore
        Crawler["Crawler (DFS)"]:::search
        Indexer["Indexer"]:::search
        SearchEngine["Search Engine"]:::search
        Tokenizer["Tokenizer"]:::search
        Crawler --> Indexer
        Indexer --> SearchEngine
        SearchEngine --> Tokenizer
    end

    %% Data Layer
    subgraph DataLayer
        DB[(PostgreSQL)]:::database
        SearchIndex["Search Indices"]:::database
        URLStorage["URL Storage"]:::database
        UserMgmt["User Management"]:::database
        DB --> SearchIndex
        DB --> URLStorage
        DB --> UserMgmt
    end

    %% Background Services
    subgraph Background
        Cron["Cron Jobs"]:::background
    end

    %% Connections between layers
    Frontend --> Server
    Server --> SearchCore
    SearchCore --> DataLayer
    Background --> SearchCore

    %% Click events for component mapping
    click Templates "https://github.com/DukeBWard/simple-search/blob/main/views/index.templ"
    click Routes "https://github.com/DukeBWard/simple-search/blob/main/routes/routes.go"
    click Auth "https://github.com/DukeBWard/simple-search/blob/main/utils/jwt.go"
    click Crawler "https://github.com/DukeBWard/simple-search/blob/main/search/crawler.go"
    click Indexer "https://github.com/DukeBWard/simple-search/blob/main/search/indexer.go"
    click SearchEngine "https://github.com/DukeBWard/simple-search/blob/main/search/engine.go"
    click Tokenizer "https://github.com/DukeBWard/simple-search/blob/main/search/tokenizer.go"
    click SearchIndex "https://github.com/DukeBWard/simple-search/blob/main/db/search_index.go"
    click URLStorage "https://github.com/DukeBWard/simple-search/blob/main/db/url.go"
    click UserMgmt "https://github.com/DukeBWard/simple-search/blob/main/db/user.go"
    click Cron "https://github.com/DukeBWard/simple-search/blob/main/utils/cron.go"
    click Fiber "https://github.com/DukeBWard/simple-search/blob/main/main.go"

    %% Legend
    subgraph Legend
        Frontend_L["Frontend"]:::frontend
        Backend_L["Backend"]:::backend
        Search_L["Search Engine"]:::search
        Database_L["Database"]:::database
        Background_L["Background Services"]:::background
    end
```
# Notes
Using htmx, tailwindcss and daisyui from CDN

# How to search
`http://127.0.0.1:port/search` body of `{"term": "term goes here"}`

# Backend
* go
* fiber
* htmx
* cron
* postgresql

# Frontend
* templ
* air for live reload: `air`
  * useful for building with templ
* daisyui
* tailwindcss

# Crawler
Uses DFS algorithms to scan and save URLs and metadata from the internet

# Indexer
Implements url indexing 

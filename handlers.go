package main

import (
	"fmt"
	"strings"
    "./cli"
	"./client"
	"./drive"
)

func listHandler(ctx cli.Context) {
    args := ctx.Args()
    gdrive := newDrive()

    gdrive.List(drive.ListFilesArgs{
        MaxFiles: args.Int64("maxFiles"),
        NameWidth: args.Int64("nameWidth"),
        Query: args.String("query"),
        SkipHeader: args.Bool("skipHeader"),
        SizeInBytes: args.Bool("sizeInBytes"),
    })
}

func downloadHandler(ctx cli.Context) {
    args := ctx.Args()
    gdrive := newDrive()

    gdrive.Download(drive.DownloadFileArgs{
        Id: args.String("id"),
        Force: args.Bool("force"),
        Stdout: args.Bool("stdout"),
        NoProgress: args.Bool("noprogress"),
    })
}

func uploadHandler(ctx cli.Context) {
    args := ctx.Args()
    gdrive := newDrive()

    gdrive.Upload(drive.UploadFileArgs{
        Path: args.String("path"),
        Name: args.String("name"),
        Parent: args.String("parent"),
        Mime: args.String("mime"),
        Recursive: args.Bool("recursive"),
        Stdin: args.Bool("stdin"),
        Share: args.Bool("share"),
    })
}

func infoHandler(ctx cli.Context) {
    args := ctx.Args()
    gdrive := newDrive()

    gdrive.Info(drive.FileInfoArgs{
        Id: args.String("id"),
        SizeInBytes: args.Bool("sizeInBytes"),
    })
}

func deleteHandler(ctx cli.Context) {
    fmt.Println("Deleting...")
}

func handler(ctx cli.Context) {
    fmt.Println("handler...")
}

func printVersion(ctx cli.Context) {
    fmt.Printf("%s v%s\n", Name, Version)
}

func printHelp(ctx cli.Context) {
    fmt.Printf("%s usage:\n\n", Name)

    for _, h := range ctx.Handlers() {
        fmt.Printf("%s %s  (%s)\n", Name, h.Pattern, h.Description)
    }
}

func printCommandHelp(ctx cli.Context) {
    handlers := ctx.FilterHandlers(ctx.Args().String("subcommand"))

    if len(handlers) == 0 {
        ExitF("Subcommand not found")
    }

    if len(handlers) > 1 {
        ExitF("More than one matching subcommand, be more specific")
    }

    handler := handlers[0]

    fmt.Printf("%s %s  (%s)\n", Name, handler.Pattern, handler.Description)
    for name, flags := range handler.Flags {
        fmt.Printf("\n%s:\n", name)
        for _, flag := range flags {
            fmt.Printf("  %s  (%s)\n", strings.Join(flag.GetPatterns(), ", "), flag.GetDescription())
        }
    }
}

// TODO: take app path as arg
func newDrive() *drive.Drive {
    oauth, err := client.NewOauthClient(ClientId, ClientSecret, DefaultTokenFilePath, authCodePrompt)
    if err != nil {
        ExitF("Failed getting oauth client: %s", err.Error())
    }

    client, err := client.NewClient(oauth)
    if err != nil {
        ExitF("Failed getting drive: %s", err.Error())
    }

    return drive.NewDrive(client)
}

func authCodePrompt(url string) func() string {
    return func() string {
        fmt.Println("Authentication needed")
        fmt.Println("Go to the following url in your browser:")
        fmt.Printf("%s\n\n", url)
        fmt.Print("Enter verification code: ")

        var code string
        if _, err := fmt.Scan(&code); err != nil {
            fmt.Printf("Failed reading code: %s", err.Error())
        }
        return code
    }
}

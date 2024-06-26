package main

import (
	"context"
	"fmt"
)

type AwsAnnouncer struct{}

// Publish the application container after building and testing it on-the-fly
func (m *AwsAnnouncer) Publish(ctx context.Context, source *Directory, repoName, repoAlias string, imageTags ...string) ([]string, error) {
	// ecrClient, err := NewECRClient(ctx)
	// if err != nil {
	// 	return "", err
	// }

	// hash, err := dirhash.HashDir(sourcePath, "", dirhash.Hash1)
	// if err != nil {
	// 	return nil, err
	// }

	tags := []string{"latest"}
	tags = append(tags, imageTags...)
	addresses := []string{}

	// ecrClient.BatchGetImage(ctx, &ecr.BatchGetImageInput{
	// 	ImageIds: []ecrtypes.ImageIdentifier{
	// 		{
	// 			ImageTag: aws.String(hash),
	// 		},
	// 	},
	// 	RepositoryName: aws.String(repoName),
	// 	RegistryId:     aws.String(registryId),
	// })

	ctr := m.Build(source)

	for _, tag := range tags {
		address, err := ctr.Publish(ctx, fmt.Sprintf("public.ecr.aws/%s/%s:%s", repoAlias, repoName, tag))
		if err != nil {
			return addresses, err
		}
		addresses = append(addresses, address)
	}

	return addresses, nil
}

// Build the application container
func (m *AwsAnnouncer) Build(source *Directory) *Container {
	build := m.BuildEnv(source).
		WithExec([]string{"go", "build", "-tags", "lambda.norpc", "-o", "dist/main"}).
		Directory("./dist")
	return dag.Container().From("public.ecr.aws/lambda/provided:al2023").
		WithDirectory("/app", build).
		WithEntrypoint([]string{"/app/main"})
}

// Return the result of running unit tests
func (m *AwsAnnouncer) Test(ctx context.Context, source *Directory) (string, error) {
	return m.BuildEnv(source).
		WithExec([]string{"go", "test"}).
		Stdout(ctx)
}

// Build a ready-to-use development environment
func (m *AwsAnnouncer) BuildEnv(source *Directory) *Container {
	goCache := dag.CacheVolume("go")
	return dag.Container().
		From("golang:1.22.3").
		WithDirectory("/src", source).
		WithMountedCache("/go/pkg/mod", goCache).
		WithWorkdir("/src").
		WithExec([]string{"go", "mod", "download"})
}

// // Returns a container that echoes whatever string argument is provided
// func (m *AwsAnnouncer) ContainerEcho(stringArg string) *Container {
// 	return dag.Container().From("alpine:latest").WithExec([]string{"echo", stringArg})
// }

// // Returns lines that match a pattern in the files of the provided Directory
// func (m *AwsAnnouncer) GrepDir(ctx context.Context, directoryArg *Directory, pattern string) (string, error) {
// 	return dag.Container().
// 		From("alpine:latest").
// 		WithMountedDirectory("/mnt", directoryArg).
// 		WithWorkdir("/mnt").
// 		WithExec([]string{"grep", "-R", pattern, "."}).
// 		Stdout(ctx)
// }

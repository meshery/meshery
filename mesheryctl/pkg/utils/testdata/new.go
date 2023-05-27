import (
    "context"
    "time"
)

func pingKubernetesCluster() {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()


    client, err := createKubernetesClient()
    if err != nil {

    }

    if err := client.Ping(ctx); err != nil {

    }


    if ctx.Err() == context.DeadlineExceeded {

    }
}
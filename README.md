# machine-status-api
RaspberryPiのGPIOを使用して、物理マシンの電源オンオフ、ProxmoxVEの仮想マシンの管理をREST API経由で行えるようにするプログラムです。

## Circuit diagram

未作成

回路写真（From Twitter）
[![Circuit Picture](./.resources/cicuitpicture.jpg)](https://twitter.com/walnuts1018/status/1628759384414367751?s=20)

## Docker Image

Buildxを利用し、arm64, amd64両方に対応

```bash
docker buildx build --platform linux/amd64,linux/arm64 -t ghcr.io/walnuts1018/machine-status-api:latest -t ghcr.io/walnuts1018/machine-status-api:<tag> . --push
```

(TODO: GitHub Actions)

## Kubernetes Manifest Sample

- [./.k8s](./.k8s)

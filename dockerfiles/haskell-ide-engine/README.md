# How to update the sourcegraph/haskell-ide-engine image

1. Open [./Dockerfile](./Dockerfile) and change the https://github.com/haskell/haskell-ide-engine commit:

```diff
   RUN git clone https://github.com/haskell/haskell-ide-engine --recursive /tmp/haskell-ide-engine \
   && cd /tmp/haskell-ide-engine \
-  && git checkout 562ac94d245e7b6ffa380eae4b02a832c397cfbb \
+  && git checkout 753deb91f9b9b39f14722315277d9fd587716bde \
   # Avoid invalidating the layers when new commits are added.
```

2. Then build the image and tag it with the first 7 characters of the commit hash (e.g. 753deb9 - this is a convention):

```
docker build -t sourcegraph/haskell-ide-engine:<NEWTAG> -f dockerfiles/haskell-ide-engine/Dockerfile .
```

3. Push the image:

```
docker push sourcegraph/haskell-ide-engine:<NEWTAG>
```


## TikToken (cl100k_base)

### Byte pair ecoding

The byte pair "aa" occurs most often, so it will be replaced by a byte that is not used in the data, such as "Z". Now there is the following data and replacement table:

```text
ZabdZabac
Z=aa
```
Then the process is repeated with byte pair "ab", replacing it with "Y":

```text
ZYdZYac
Y=ab
Z=aa
```

https://en.wikipedia.org/wiki/Byte_pair_encoding


### clk100k

https://openaipublic.blob.core.windows.net/encodings/cl100k_base.tiktoken



## Embeddings (text-embedding-ada-002)


https://platform.openai.com/docs/api-reference/embeddings/object


``
curl https://api.openai.com/v1/embeddings -H "Content-Type: application/json" -H "Authorization: Bearer $OPENAI_API_KEY" -d '{"input": "Your text string goes here","model": "text-embedding-ada-002"}'
``


### Local Database

```shell
docker pull ankane/docker-pgvector
```



```shell

find . -type f -name "*.go" -print0 | xargs -0 cat

```
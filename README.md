This is an implementation of the issue here: [https://gist.github.com/ehmo/e736c827ca73d84581d812b3a27bb132]
```
Below is a list of links leading to an image, read this list of images and find 3 most prevalent colors in the RGB scheme in hexadecimal format (#000000 - #FFFFFF) in each image, and write the result into a CSV file in a form of url,color,color,color.

Please focus on speed and resources. The solution should be able to handle input files with more than a billion URLs, using limited resources (e.g. 1 CPU, 512MB RAM). Keep in mind that there is no limit on the execution time, but make sure you are utilizing the provided resources as much as possible at any time during the program execution.

Answer should be posted in a git repo.
```

## Approach

The URLs in the input file are batched to allow manageable processing with limited resources. The batch size - number of urls in a batch size should be configurable. Each batch is saved in a temporary file as `<input file name>_p{n}.txt`, where `n` is 1,2,3 ... N, e.g. `input_p1.txt`. The temporary partition files are removed after processing. The size of the batch (default 100) is passed as an argument to the program: `--pSize=200`

Concurrent processing of the batches is also supported to improve speed. The concurrency is achieved with worker goroutines that process the partition files. The number of workers is configurable with a program parameter: e.g., `--workers=3` (default 2).

The input file is specified with parameter `inputFile` (default  `input.txt`). 

The built executable is run as:
```
pexChallenge -inputFile='input.txt' -pSize=200 -workers=5
```

## Memory Usage
The internal result of the processing is stored in type `imageResult`. It has a `Colors` field which is an array of size 3 (for the three most prevalent colors of the URL). Each element is a 3 byte/unsigned 8 bit array for the R, G, B components of the color. A total of 9 bytes each used for the result of each URL.

An alternative option would be to use the values in hexadecimal format. However, this would require a total of 18 bytes for the result of each URL - 2 (bytes for each color component) X 3 (number of color components) X 3 (number of colors).


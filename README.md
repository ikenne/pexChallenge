This is an implementation of the issue here: [https://gist.github.com/ehmo/e736c827ca73d84581d812b3a27bb132]
```
Below is a list of links leading to an image, read this list of images and find 3 most
prevalent colors in the RGB scheme in hexadecimal format (#000000 - #FFFFFF) in 
each image, and write the result into a CSV file in a form of url,color,color,color.

Please focus on speed and resources. The solution should be able to handle input files
with more than a billion URLs, using limited resources (e.g. 1 CPU, 512MB RAM).
Keep in mind that there is no limit on the execution time, but make sure you are
utilizing the provided resources as much as possible at any time during the program 
execution.

Answer should be posted in a git repo.
```

## Approach

The URLs in the input file are batched to allow manageable processing with limited resources. The batch size - number of urls in a batch - should be configurable. Each batch is saved in a temporary partition file as `<input file name>_p{n}.txt`, where `n` is 1,2,3, ... N, e.g., `input_p1.txt`. The partition files are removed after processing. The size of the batch is passed as an argument to the program: e.g.,`--pSize=200` (default 100).

Concurrent processing of the batches is also supported to improve speed. The concurrency is achieved by processing the partition files in worker goroutines and using buffered channels. The number of workers is configurable with a program parameter: e.g., `--workers=3` (default 2). Buffered channels are used for queuing and coordinating processing and writing the results between the goroutines.

The input file is specified with parameter `inputFile` (default  `input.txt`). 

The built executable is run as:
```
pexChallenge -inputFile='input.txt' -pSize=200 -workers=3
```

The CSV output is saved as `output.txt`

## Memory Usage
The internal result of the processing is stored in type `imageResult` (model/imageResult.go). It has a `Colors` field which is an array of size 3 (for the three most prevalent colors of the URL). Each array item is a 3 byte (unsigned 8 bit int) array for the R, G, B color components, thus a total of 9 bytes each used for the result of each URL.

An alternative option would be to store the values in hexadecimal string format. However, this would require a total of 18 bytes for the result of each URL - 2 (bytes for each color component in hexadecimal format) X 3 (number of color components) X 3 (number of colors).

Therefore the first option uses about 50% memory in comparison to the alternative option.

## Execution times
The following is a sample of local execution times with different configurations:
```
pSize=200, workers=1: Elapsed time 108.596055 seconds
pSize=200, workers=2: Elapsed time 78.3090502 seconds
pSize=200, workers=3: Elapsed time 63.5690846 seconds
pSize=200, workers=4: Elapsed time 67.0927973 seconds
pSize=200, workers=5: Elapsed time 58.622737 seconds
```


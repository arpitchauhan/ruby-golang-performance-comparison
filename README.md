### Comparing Golang and Ruby CSV Validation Performance

This project has various CSV validation scripts written in the two languages, and you can run them with the included CSV file (which has a million rows) to compare their performance.

A script `benchmark.rb` has been included, and it can run all the validation scripts at once and report the times taken by them to you.  

To run the benchmark, clone this repository and make sure that you have working Golang and Ruby installations.

And then, run this from the cloned directory in a shell:

```
ruby benchmark.rb
```

If you are unwilling to run this yourself, you can also see the [sample output](sample_output.txt) that I got from running the benchmark on my machine.
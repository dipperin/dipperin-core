### Building the docs on your machine

Here are the quick steps to achieve this on a local machine without depending on ReadTheDocs, starting from the main fabric directory. Note: you may need to adjust depending on your OS. 

```
sudo pip install Sphinx
sudo pip install sphinx_rtd_theme
sudo pip install recommonmark==0.4.0
sudo pip install sphinx-markdown-tables
cd ~/go/src/github.com/dipperin/dipperin-core/docs # Be in this directory. Makefile sits there.
make html
```
This will generate all the html files in ```docs/build/html``` which you can then start browsing locally using your browser. ```index.html``` is the home page.

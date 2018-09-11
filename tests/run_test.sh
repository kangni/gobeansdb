#!/bin/bash

virtualenv venv
source venv/bin/activate
venv/bin/python venv/bin/pip install Cython>=0.20 -i https://pypi.douban.com/simple
venv/bin/python venv/bin/pip install -r tests/pip-req.txt -i https://pypi.douban.com/simple
venv/bin/python venv/bin/nosetests --with-xunit --xunit-file=unittest.xml
deactivate

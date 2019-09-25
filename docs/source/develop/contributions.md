
# Contributions

The Dipperin project eagerly accepts contributions from the community.We welcome contributions to Dipperin in many forms.
## Working Together
When contributing or otherwise participating, please:         
* Be friendly and welcoming   
* Be patient   
* Be thoughtful   
* Be respectful   
* Be charitable   
* Avoid destructive behavior   

Excerpted from the [Go conduct document](https://golang.org/conduct).

## Ways to contribute
### Getting help
If you are looking for something to work on, or need some expert assistance in debugging a problem or working out a fix to an issue, our community is always eager to help. We hang out on mail `report@dipperin.com`. Questions are in fact a great way to help improve the project as they highlight where our documentation could be clearer.

### Reporting Bugs 
When you encounter a bug, please open an issue on the corresponding repository. Start the issue title with the repository/sub-repository name, like ```repository_name: issue name```. We have provided a issue templates for bug report:Bug report template. If you can abide by this template, this will help us fix the bug more efficiently.

### Suggesting Enhancements
If the scope of the enhancement is small, open an issue. If it is large, such as suggesting a new repository, sub-repository, or interface refactoring, then please @Dipperin-Project on an issue,we will pay more attention on you suggestion. 

### Your First Code Contribution
If you are a new contributor, thank you! Before your first merge, you will need to be added to the [CONTRIBUTORS](https://github.com/dipperin/dipperin-core/blob/dev/CONTRIBUTORS) files. Open a pull request adding yourself to these files. All Dipperin code follows the LGPL license in the license document. We prefer that code contributions do not come with additional licensing. For exceptions, added code must also follow a LGPL license.

### Code Contribution
If it is possible to split a large pull request into two or more smaller pull requests, please try to do so. 
Pull requests should include tests for any new code before merging. It is ok to start a pull request on partially implemented code to get feedback, and see if your approach to a problem is sound. 
You don't need to have tests, or even have code that compiles to open a pull request, although both will be needed before merge. When tests use magic numbers, please include a comment explaining the source of the number.    
Commit messages also follow some rules. They are best explained at the official [Go](https://golang.org/) "Contributing guidelines" document:
[golang.org/doc/contribute.html](https://golang.org/doc/contribute.html#commit_changes)

For example:   

```
Dipperin-core: add support for consensus
	
This change list adds support for consensus.
Previously, the Dipperin-core package was consensus slowly,sometimes leading to
a panic later on in the program execution.
Improve consensus efficiency and add some tests.
	
Fixes  Dipperin/Dipperin-core/core#20.
```
If the change list modifies multiple packages at the same time, include them in the commit message:   

```
Dipperin-core/core,Dipperin-core/core/Dipperin: implement wrapping of Go interfaces

bla-bla

Fixes Dipperin/Dipperin-core/core#40.
```
Please always format your code with [goimports](https://godoc.org/golang.org/x/tools/cmd/goimports). Best is to have it invoked as a hook when you save your .go files.

Files in the Dipperin repository don't list author names, both to avoid clutter and to avoid having to keep the lists up to date. Instead, your name will appear in the change log and in the [CONTRIBUTORS](https://github.com/dipperin/dipperin-core/blob/dev/CONTRIBUTORS) files.

New files that you contribute should use the standard copyright header:

```
// Copyright 2019, Keychain Foundation Ltd.
// This file is part of the dipperin-core library.
//
// The dipperin-core library is free software: you can redistribute
// it and/or modify it under the terms of the GNU Lesser General Public License
// as published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
//
// The Dipperin-core library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.
```

Files in the repository are copyright the year they are added. Do not update the copyright year on files that you change.

## Code Review
We follow the convention of requiring at least 1 reviewer to say LGTM(looks good to me) before a merge. When code is tricky or controversial, submitters and reviewers can request additional review from others and more LGTMs before merge. You can ask for more review by saying PTAL(please take another look) in a comment in a pull request. You can follow a PTAL with one or more @someone to get the attention of particular people. If you don't know who to ask, and aren't getting enough review after saying PTAL, then PTAL @Dipperin-Project will get more attention. Also note that you do not have to be the pull request submitter to request additional review.

## Style
We use [Go style](https://github.com/golang/go/wiki/CodeReviewComments).

## What Can I Do to Help?
If you are looking for some way to help the Dipperin project, there are good places to start, depending on what you are comfortable with.   
   You can search for open issues in need of resolution.   
   You can improve documentation, or improve examples.   
   You can add and improve tests.   
   You can improve performance, either by improving accuracy, speed, or both.   
   You can suggest and implement new features that you think belong in Dipperin.    

**********************

This "Contributing" guide has been extracted from the [Gonum](https://www.gonum.org/) project. Its guide is [here](https://github.com/gonum/license/blob/master/CONTRIBUTING.md).

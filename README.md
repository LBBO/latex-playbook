# Latex playbook
A LaTeX project to create nice looking theater playbooks, as used by the german ensemble [Der ChaosTRUB](http://chaostrub.de).

## Motivation
We often face the problem that our plays are formatted in a way that we don't particularly care for. However, reformatting everything in a WYSIWYG editor is very time-consuming and not exactly enjoyable.

We discovered that there are usually simple ways to add / remove text automatically to the places we would like to change, but text usually doesn't affect the formatting of documents. Well, this is where LaTeX comes in! You can automatically (via a script, RegExp or something similar) add the commands to the playbook and then just paste it into the `tex` file - voilà, you're done!

## Usage
### Setup LaTeX
As this is a 100%-LaTeX project, you will need to have that installed before continuing. For more information, please check [this website](https://www.latex-project.org/get/).

### Setup characters
The LaTeX file will automatically count the characters' appearances for you, but it needs some help getting started.

First of all, find the title page (i.e. by searching for `\begin{titlepage}`). This page is going to require some setup.

#### Create a counter for each character
There is a block where counters are created for each and every character. Replace the existing lines with one line per character that appears in your play:

```tex
\newcounter{Name of character}
```

**ATTENTION:** The character's name is case-sensitive here! It has to be spelled **exactly** as it will be throughout the rest of the play.

**ATTENTION:** The LaTeX file must be compiled twice to correctly output the couters' results! The first compilation after a change will only result in question marks instead of correct numbers.

#### Fill out the table of characters
Below the counters, there is a table that states the characters' names, how often they appear, what they are like, and who they are played by. Most of this information will have to be filled out manually, except for the number of appearances: insert `\ref{Name of character}` instead.

Sometimes, the characters lines start with a shorter version of their name. To clarify how they will be referred to during the play, we tend to print that part of their name in bold font.

**ATTENTION:** in the `\ref` command, the name must be spelled **exactly** as it was spelled in the `\newcounter` command!

```tex
% Copy this line once for every character and replace the example data
\textbf{ALICE} SMITH & \ref{Alice} & Mysterious lady & Jane Doe
```

#### Fill out the table of scenes and appearances
This table gives the actors and actrisses a quick overview of who appears in which scene and which pages the scene is on. Create one column per character and one row per scene. You'll need to fill in the page numbers manually (for now) and then add information about who is needed during every scene.

Our convention is to make the table cell dark grey (we added a color called `TableColorAppearance` for this) if the character says anything in a scene, light gray (we added a color called `TableColorSemiAppearance` for this) with a "+" if the character doesn't have any lines but technically is still present, and blank if the character is absent during the entire scene.

This table's code should look somewhat like this:

```tex
Scene & Pages & Character 1 & Character 2 & ... & Character n \\ \hline
1 & 1-2 &  &  & ... & \cellcolor{TableColorAppearance} \\ \hline % In this scene, Character n might be having a monologue
2 & 2-5 & \cellcolor{TableColorAppearance} & \cellcolor{TableColorAppearance} & ... & \cellcolor{TableColorSemiAppearance} \\ \hline % Now Character 1 and Character 2 might be discussing something while Character n is just observing from the background
...
k & i-j &  &  & ... & \cellcolor{TableColorAppearance} \\ \hline % Perhaps Character n ends the play with another monologue
```

## Insert the play
Now comes the main part! In order to correctly format the play, please use the following commands:

* `\scene{Scene title}{Character 1, Character 2 (+ Character 3, Character 4)}` - Insert a scene's title and, below the title, note which characters appear in this scene. We tend to even add characters that don't say anything but are present nonetheless in brackets.

* `\character{Name} Lorem ipsum dolor sit amet...` - Begin a character's line. **ATTENTION** The character's name must be spelled **exactly** as in the previous commands (including lower-case and upper-case letters)!

* `\InlineStageDirection{Stage direction}` - Add a small stage direciton (i.e. "whispers", "grabs the box") in the middle of a line. Please note the curly braces around the stage direction!

* `\BlockStageDirection Stage direction` - Add a large stage direction (i.e. multiple characters entering the stage and interacting with each other without saying anything) that is important enough to be "its own line". Please note the absence of curly braces around the actual stage direction!


## Customize the playbook
Time to get creative! You might want to customize the very first page (which is currently directed towards strangers, asking them to return the playbook in case it is lost by a cast member and then found by somebody) by adding your group's contact information.

You could also add some metadata about your play above or below the table of your characters (or anywhere else) and just customize the design until you're happy with it!

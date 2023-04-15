Output demo.gif

Set FontSize 14
Set Width 1200
Set Height 600

Hide
Type "source mop.bash"
Enter
Show

Type "# Welcome to mop!"
Sleep 1s
Enter
Sleep 1s

Type "# Before you begin, make sure to create a GitHub Issue that defines the operation. We will use Issue #1 in this demo."
Sleep 1s
Enter
Type "gh issue view 1"
Sleep 1s
Enter
Sleep 7s

Ctrl+L
Sleep 1s

Type "# To start an operation, execute the 'mop-start' command."
Sleep 1s
Enter
Type "mop-start 1"
Sleep 1s
Enter
Sleep 5s

Type "# Press 'Ctrl + t' to auto-complete the first command."
Sleep 1s
Enter
Ctrl+T
Sleep 2s
Enter
Sleep 2s

Type "# Run the 'next' command to proceed to the following command in the sequence."
Sleep 1s
Enter
Type "next"
Enter
Sleep 1s

Type "# Press 'Ctrl + t' again to auto-complete the second command."
Sleep 1s
Enter
Ctrl+T
Sleep 2s
Enter
Sleep 2s

Type "# Run 'upload' command to post the execution results in the comments of the GitHub Issue."
Sleep 1s
Enter
Type "upload"
Sleep 1s
Enter
Sleep 3s

Type "# Verify that the results have been successfully uploaded."
Sleep 1s
Enter
Type "gh issue view 1 --comments"
Enter
Sleep 7s

Ctrl+L
Sleep 1s

Type "# Additionally, press 'Ctrl + o' to use pre-defined useful commands using fuzzy search like fzf."
Sleep 1s
Enter
Ctrl+O
Sleep 2s
Up
Sleep 2s
Up
Sleep 2s
Enter
Sleep 1s
Enter

Sleep 5s

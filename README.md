# cred-gen
A simple cli app to easily get Salesforce sessionIds

Fulfill the .env.example in a new `*.sh` file with execute permission and source it in your shell config file

Example:
1) cd ~ 
2) `touch sfenvs.sh` 
3) Export all the variables on .env.example inside the new sh file
4) Give execution permission with: `chmod +x sfenvs.sh` 
5) Source it in zsh config file, for example: `echo "$HOME/.sfenvs.sh" >> .zshrc`

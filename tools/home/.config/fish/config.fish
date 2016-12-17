cd "$CODE"

# Intialize git config.
git config --local include.path ../tools/git/config

# Initialize git-subrepo.
bass source "$CODE"/vendor/github.com/ingydotnet/git-subrepo/.rc

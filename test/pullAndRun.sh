## declare an array variable
declare -a repos=(
"https://github.com/wloop/wlib-malloc"
"https://github.com/wloop/wlib-timer"
"https://github.com/wloop/wlib-tlsf"
"https://github.com/wloop/wlib-memory"
"https://github.com/wloop/wlib-json"
"https://github.com/wloop/wlib-queue"
"https://github.com/wloop/wlib-fsm"
"https://github.com/wloop/wlib-tmp"
)

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}"  )" >/dev/null && pwd  )"

rm -rf "$DIR/repo-test"
mkdir "$DIR/repo-test"
cd "$DIR/repo-test"

## now loop through the above array
for i in "${repos[@]}"
do
    # clone the repo
    git clone --depth 1 "$i"

    # get repo name
    basename=$(basename "$i")
    filename=${basename%.*}

    # cd inro the repo
    cd "$DIR/repo-test/$filename"

    # execute wio and test
    wio clean --hard
    wio update
    wio install
    wio build --all

    cd ../
done

rm -rf "$DIR"/repo-test

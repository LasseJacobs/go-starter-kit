GIT_DIGEST=latest

rm -r build/ || echo 'no directory' 
mkdir build/
cp *.yaml build/
sed -i '' 's/\$GIT_TAG/$GIT_DIGEST/g' build/*.yaml

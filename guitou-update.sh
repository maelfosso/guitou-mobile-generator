#!/bin/bash


check_error() {
  if (( $? )); then
    exit $1
  fi
}

printf "\n"

PACKAGE=${1:?"Error. Must set the package ID"}
PROJECT_NAME=${2:?"Error. Must set the project NAME"}
REPO="/tmp/guitou-mobile"
DIR="/tmp/guitou-$PACKAGE"

printf "\t%s\t%s\t%s\n"  "[Guitou Mobile Generator] Let's start. " ${DIR} $PACKAGE

# Delete folder project if it exists
# [ -d "${DIR}" ] \
#   && { rm -rf ${DIR}; printf "\t%s\n"  "[Guitou Mobile Generator] Folder deleted if exists $DIR"; }

# [ -d "${DIR}" ] \
#   && { cd $DIR; echo "[Guitou Mobile Generator] Changing Directory"; }
#   || { printf "\t" "[Guitou Mobile Generator] Directory doesn't exists"; exit 1; }

# Change the app name
ANDROID_MANIFEST="$DIR/android/app/src/main/AndroidManifest.xml"
printf "\t%s\t%s\n" "[Guitou Mobile Generator] Changing the App Name" $DIR
# cat $ANDROID_MANIFEST

sed -i "s/Guitou/$PROJECT_NAME/g" $ANDROID_MANIFEST
check_error 2
printf "\t%s\n"  "[Guitou Mobile Generator] App name changed"


BASE="$DIR/android/app/src/main/kotlin"
cd $BASE

NEW_PKG_PATH="cm/guitou/mobile/$PACKAGE/"
TMP_PKG_PATH="tmp/"
OLD_PKG_PATH="cm/guitou/android/xorms/"
printf "\t%s\n"  $NEW_PKG_PATH $OLD_PKG_PATH

cp -r $OLD_PKG_PATH/ $TMP_PKG_PATH
rm $OLD_PKG_PATH/*
rmdir -p $OLD_PKG_PATH

mkdir -p $NEW_PKG_PATH
cp $TMP_PKG_PATH/* $NEW_PKG_PATH
rm -rf $TMP_PKG_PATH

cd $DIR
grep -rl "cm.guitou.android.xorms" android/* \
  | xargs sed -i "s/cm.guitou.android.xorms/cm.guitou.mobile.$PACKAGE/g"

check_error 3
printf "\t%s\n"  "[Guitou Mobile Generator] Package name changed"

# Copy the project json data into assets

printf "\t%s\n"  "[Guitou Mobile Generator] Deleting Assets"
if [ -r "$DIR/assets/project.json" ] 
then
  # cat "$DIR/assets/project.json"
  rm $DIR/assets/project.json
  # cp ./assets/project.json $DIR/assets

  check_error 4
  printf "\t%s\n"  "[Guitou Mobile Generator] Asset deleted"
fi
# rm $DIR/assets/project.json && cp ./assets/project.json $DIR/assets/

exit 0
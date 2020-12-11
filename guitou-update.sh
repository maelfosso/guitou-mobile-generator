#!/bin/bash


check_error() {
  if (( $? )); then
    exit $1
  fi
}

PACKAGE=${1:?"Error. Must set the package ID"}
REPO="/tmp/guitou-mobile"
DIR="/tmp/guitou-$PACKAGE"

echo "[Guitou Mobile Generator] Let's start. " ${DIR}

# Delete folder project if it exists
[ -d "${DIR}" ] \
  && { rm -rf ${DIR}; echo "[Guitou Mobile Generator] Folder deleted if exists $DIR"; }


# Change the working directory
cp -r $REPO $DIR
check_error 1
echo "[Guitou Mobile Generator] Working directory changed"

# Change the app name
ANDROID_MANIFEST="$DIR/android/app/src/main/AndroidManifest.xml"
sed "s/Guitou/$PACKAGE/g" $ANDROID_MANIFEST
check_error 2
echo "[Guitou Mobile Generator] App name changed"

# Change the icon (load in the Splash)


# Change the package to cm.guitou.mobile.${PACKAGE}

BASE="$DIR/android/app/src/main/kotlin"
cd $BASE

NEW_PKG_PATH="cm/guitou/mobile/$PACKAGE/"
TMP_PKG_PATH="tmp/"
OLD_PKG_PATH="cm/guitou/android/xorms/"
echo $NEW_PKG_PATH $OLD_PKG_PATH

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
echo "[Guitou Mobile Generator] Package name changed"

# Copy the project json data into assets

rm $DIR/assets/project.json && cp ./assets/project.json $DIR/assets/
check_error 4
echo "[Guitou Mobile Generator] Asset copied"


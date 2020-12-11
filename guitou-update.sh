#!/bin/bash


check_error() {
  if (( $? )); then
    exit 1
  else

  fi
}

REPO="/tmp/guitou-mobile"
DIR="/tmp/guitou-$PROJECT_NAME"

# Change the working directory
mv $REPO $DIR
check_error()

# Change the app name
ANDROID_MANIFEST= "$DIR/android/app/src/main/AndroidManifest.xml"
sed "s/Guitou/$PROJECT_NAME/g" $ANDROID_MANIFEST
check_error()

# Change the icon (load in the Splash)


# Change the package to cm.guitou.mobile.${PROJECT_NAME}

BASE="$DIR/android/app/src/main/kotlin"
NEW_PKG_PATH="$BASE/cm/guitou/mobile/$PROJECT_NAME"
OLD_PKG_PATH="$BASE/cm/guitou/android/xorms"

mv $OLD_PKG_PATH $NEW_PKG_PATH
grep -rl "cm.guitou.android.xorms" android/* \ 
  | xargs sed "s/cm.guitou.android.xorms/cm.guitou.mobile.$PROJECT_NAME/g"

check_error()

# Copy the project json data into assets

rm $DIR/assets/project.json && cp ./assets/project.json $DIR/assets/
check_error()


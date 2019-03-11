import QtQuick 2.9
import QtQuick.Controls 2.2
import QtQuick.Controls.Material 2.2

Item {
  id: root
  height: 480
  width: 114

  property bool faderEnabled
  property bool muted
  property int faderValue
  property string label
  property real meterValue

  property string oscAddress

  signal faderMoved(string addr, real pos)
  signal muteToggled(string addr, bool checked)

  Rectangle {
    anchors.horizontalCenter: parent.horizontalCenter
    y: 204
    width: root.width*0.5175
    height: 1
    color: "grey"
  }

  Slider {
    id: fader
    anchors.horizontalCenter: parent.horizontalCenter
    y: 113
    width: parent.width*0.6526
    height: 329
    to: 100
    orientation: Qt.Vertical
    value: parent.faderValue
    onMoved: parent.faderMoved(parent.oscAddress, position)
    enabled: parent.faderEnabled
    Material.accent: enabled ? defaultAccent : defaultDisabled
  }

  Button {
      id: mute
      anchors.horizontalCenter: parent.horizontalCenter
      anchors.topMargin: 10
      anchors.top: parent.top
      width: parent.width*0.80
      height: 60
      text: qsTr("mute")
      scale: 1
      checkable: true
      checked: parent.muted
      onToggled: parent.muteToggled(parent.oscAddress, checked)
  }

  Text {
      anchors.horizontalCenter: parent.horizontalCenter
      anchors.top: fader.bottom
      anchors.topMargin: 5
      width: parent.width
      wrapMode: Text.Wrap
      text: parent.label
      font.pixelSize: 14
      horizontalAlignment: TextInput.AlignHCenter
  }

  ProgressBar {
    value: parent.meterValue
    width: parent.width*0.53509
    anchors.horizontalCenter: parent.horizontalCenter
    anchors.top: mute.bottom
    anchors.topMargin: 15
  }
}

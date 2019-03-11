import QtQuick 2.9
import QtQuick.Controls 2.2
import QtQuick.Controls.Material 2.2
import QtQuick.Layouts 1.2

Item {
    id: root
    width: 800
    height: 480

    property var defaultAccent: Material.Pink
    property var defaultDisabled: Material.Grey
    property var configuration: JSON.parse(controllerConfig)

    Material.accent: defaultAccent

    Item {
      id: controller_layer
      width: 800
      height: 480
      layer.enabled: true
      visible: !QmlRoot.busy


      RowLayout {
        spacing: 0

        Repeater {
          model: configuration.channelStrips.length

          ChannelStrip {
            property var config: configuration.channelStrips[index]
            oscAddress: config.oscAddress

            Layout.alignment: Qt.AlignTop
            width: 684 / configuration.channelStrips.length
            faderEnabled: !lockButton.checked

            onFaderMoved: QmlRoot.sendFaderValue(addr, pos)
            onMuteToggled: QmlRoot.sendMute(addr, checked)

            Component.onCompleted: QmlRoot.registerChannelStrip(config.oscAddress, this)
          }
        }
      }

      Item {
        id: buttonArea
        anchors.topMargin: 0
        anchors.top: parent.top
        anchors.right: parent.right
        width: 116
        height: parent.height
        visible: true

        Rectangle {
          anchors.fill: parent
          color: "grey"
        }

        Button {
          id: lockButton
          text: qsTr("Lock\nFaders")
          anchors.top: parent.top
          anchors.topMargin: 10
          anchors.horizontalCenter: parent.horizontalCenter
          width: 71
          height: 60
          checkable: true
        }

        Button {
            id: recallbutton
            anchors.horizontalCenter: parent.horizontalCenter
            anchors.bottomMargin: 10
            anchors.bottom: parent.bottom

            width: 71
            height: 60
            visible: configuration.recallButton.enabled
            text: configuration.recallButton.label
            enabled: !lockButton.checked
            onPressAndHold: confirmRecallDialog.open()
        }
      }
    }

    Dialog {
      id: confirmRecallDialog
      modal: true
      title: qsTr("Confirm %1").arg(configuration.recallButton.label)

      x: (parent.width - width) / 2
      y: (parent.height - height) / 2

      standardButtons: Dialog.Ok | Dialog.Cancel
      onAccepted: QmlRoot.recallClicked(configuration.recallButton.sceneNumber)
      Label {
        text: qsTr("Are you sure you want to do a %1 ?").arg(configuration.recallButton.label)
      }
    }

    Item {
      id: busy_layer
      visible: QmlRoot.busy

      BusyIndicator {
        id: busyindicator
        layer.enabled: true
        width: 800
        height: 480
        running: QmlRoot.busy
      }

      Text {
        width: 800
        height: 480
        horizontalAlignment: TextInput.AlignHCenter
        verticalAlignment: TextInput.AlignVCenter
        font.pixelSize: 30
        text: qsTr("Connecting...")
      }
    }

    Drawer {
        id: drawer
        width: parent.width*0.20
        height: parent.height
        edge: Qt.RightEdge

        Button {
          text: qsTr("Shutdown")
          width: 100
          height: 50
          anchors.topMargin: 10
          anchors.top: parent.top
          anchors.horizontalCenter: parent.horizontalCenter
          onPressAndHold: confirmShutdownDialog.open()
        }

        Slider {
          id: brightness
          anchors.centerIn: parent
          y: 118
          width: 69
          height: 200
          to: 250
          from: 15
          orientation: Qt.Vertical
          value: QmlRoot.brightness
          onMoved: QmlRoot.changeBrightness(value)
        }

        Text {
          text: qsTr("Display\nbrightness")
          anchors.horizontalCenter: parent.horizontalCenter
          horizontalAlignment: TextInput.AlignHCenter
          verticalAlignment: TextInput.AlignVCenter
          anchors.top: brightness.bottom
          anchors.topMargin: 10
          font.pixelSize: 12
        }

        Text {
        anchors.horizontalCenter: parent.horizontalCenter
        anchors.bottom: parent.bottom
        anchors.bottomMargin: 10
        text: QmlRoot.ipaddress
        font.pixelSize: 12
        }
    }

    Dialog {
      id: confirmShutdownDialog
      modal: true
      title: qsTr("Confirm Shutdown")

      implicitWidth: 250
      implicitHeight: 200

      x: (parent.width - width) / 2
      y: (parent.height - height) / 2

      standardButtons: Dialog.Ok | Dialog.Cancel
      onAccepted: QmlRoot.shutdown(restartSelector.checked)

      ButtonGroup {
        buttons: shutdownButtons.children
      }

      Item {
        id: shutdownButtons

        Button {
          id: restartSelector
          text: qsTr("Restart")
          checked: true
          anchors.centerIn: parent.center
          width: 150
          height: 40
          checkable: true
        }

        Button {
          id: shutdownSelector
          text: qsTr("Shutdown")
          checked: false
          anchors.top: restartSelector.bottom
          width: 150
          height: 40
          checkable: true
        }
      }
    }
}

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

      Item {
        id: chStrip0
        anchors.topMargin: 0
        anchors.top: parent.top
        anchors.left: parent.left
        width: 114
        visible: configuration.channelStrips[0].enabled

        Rectangle {
          anchors.horizontalCenter: parent.horizontalCenter
          y: 209
          width: 59
          height: 1
          color: "grey"
        }

        Slider {
          id: ch0
          anchors.horizontalCenter: parent.horizontalCenter
          y: 118
          width: 69
          height: 329
          to: 100
          orientation: Qt.Vertical
          value: strip0.faderValue
          enabled: !lockButton.checked
          onMoved: strip0.fadermoved(position)
          Material.accent: enabled ? defaultAccent : defaultDisabled
        }

        Button {
            id: mute0
            anchors.horizontalCenter: parent.horizontalCenter
            anchors.topMargin: 10
            anchors.top: parent.top
            width: 71
            height: 60
            text: qsTr("mute")
            scale: 1
            checkable: true
            checked: strip0.muted
            onToggled: strip0.muteclicked(checked)
        }

        Text {
            anchors.horizontalCenter: parent.horizontalCenter
            anchors.top: ch0.bottom
            anchors.topMargin: 10
            text: strip0.label
            font.pixelSize: 14
            horizontalAlignment: TextInput.AlignHCenter
        }

        ProgressBar {
          value: strip0.meterValue
          width: 61
          anchors.horizontalCenter: parent.horizontalCenter
          anchors.top: mute0.bottom
          anchors.topMargin: 15
        }
      }

      Item {
        id: chStrip1
        anchors.topMargin: 0
        anchors.top: parent.top
        anchors.left: chStrip0.right
        width: 114
        visible: configuration.channelStrips[1].enabled

        Rectangle {
          anchors.horizontalCenter: parent.horizontalCenter
          y: 209
          width: 59
          height: 1
          color: "grey"
        }

        Slider {
          id: ch1
          anchors.horizontalCenter: parent.horizontalCenter
          y: 118
          width: 69
          height: 329
          to: 100
          orientation: Qt.Vertical
          value: strip1.faderValue
          enabled: !lockButton.checked
          onMoved: strip1.fadermoved(position)
          Material.accent: enabled ? defaultAccent : defaultDisabled
        }

        Button {
            id: mute1
            anchors.horizontalCenter: parent.horizontalCenter
            anchors.topMargin: 10
            anchors.top: parent.top
            width: 71
            height: 60
            text: qsTr("mute")
            scale: 1
            checkable: true
            checked: strip1.muted
            onToggled: strip1.muteclicked(checked)
        }

        Text {
            anchors.horizontalCenter: parent.horizontalCenter
            anchors.top: ch1.bottom
            anchors.topMargin: 10
            text: strip1.label
            font.pixelSize: 14
            horizontalAlignment: TextInput.AlignHCenter
        }

        ProgressBar {
          value: strip1.meterValue
          width: 61
          anchors.horizontalCenter: parent.horizontalCenter
          anchors.top: mute1.bottom
          anchors.topMargin: 15
        }
      }

      Item {
        id: chStrip2
        anchors.topMargin: 0
        anchors.top: parent.top
        anchors.left: chStrip1.right
        width: 114
        visible: configuration.channelStrips[2].enabled

        Rectangle {
          anchors.horizontalCenter: parent.horizontalCenter
          y: 209
          width: 59
          height: 1
          color: "grey"
        }

        Slider {
          id: ch2
          anchors.horizontalCenter: parent.horizontalCenter
          y: 118
          width: 69
          height: 329
          to: 100
          orientation: Qt.Vertical
          value: strip2.faderValue
          enabled: !lockButton.checked
          onMoved: strip2.fadermoved(position)
          Material.accent: enabled ? defaultAccent : defaultDisabled
        }

        Button {
            id: mute2
            anchors.horizontalCenter: parent.horizontalCenter
            anchors.topMargin: 10
            anchors.top: parent.top
            width: 71
            height: 60
            text: qsTr("mute")
            scale: 1
            checkable: true
            checked: strip2.muted
            onToggled: strip2.muteclicked(checked)
        }

        Text {
            anchors.horizontalCenter: parent.horizontalCenter
            anchors.top: ch2.bottom
            anchors.topMargin: 10
            text: strip2.label
            font.pixelSize: 14
            horizontalAlignment: TextInput.AlignHCenter
        }

        ProgressBar {
          value: strip2.meterValue
          width: 61
          anchors.horizontalCenter: parent.horizontalCenter
          anchors.top: mute2.bottom
          anchors.topMargin: 15
        }
      }

      Item {
        id: chStrip3
        anchors.topMargin: 0
        anchors.top: parent.top
        anchors.left: chStrip2.right
        width: 114
        visible: configuration.channelStrips[3].enabled

        Rectangle {
          anchors.horizontalCenter: parent.horizontalCenter
          y: 209
          width: 59
          height: 1
          color: "grey"
        }

        Slider {
          id: ch3
          anchors.horizontalCenter: parent.horizontalCenter
          y: 118
          width: 69
          height: 329
          to: 100
          orientation: Qt.Vertical
          value: strip3.faderValue
          enabled: !lockButton.checked
          onMoved: strip3.fadermoved(position)
          Material.accent: enabled ? defaultAccent : defaultDisabled
        }

        Button {
            id: mute3
            anchors.horizontalCenter: parent.horizontalCenter
            anchors.topMargin: 10
            anchors.top: parent.top
            width: 71
            height: 60
            text: qsTr("mute")
            scale: 1
            checkable: true
            checked: strip3.muted
            onToggled: strip3.muteclicked(checked)
        }

        Text {
            anchors.horizontalCenter: parent.horizontalCenter
            anchors.top: ch3.bottom
            anchors.topMargin: 10
            text: strip3.label
            font.pixelSize: 14
            horizontalAlignment: TextInput.AlignHCenter
        }

        ProgressBar {
          value: strip3.meterValue
          width: 61
          anchors.horizontalCenter: parent.horizontalCenter
          anchors.top: mute3.bottom
          anchors.topMargin: 15
        }
      }

      Item {
        id: chStrip4
        anchors.topMargin: 0
        anchors.top: parent.top
        anchors.left: chStrip3.right
        width: 114
        visible: configuration.channelStrips[4].enabled

        Rectangle {
          anchors.horizontalCenter: parent.horizontalCenter
          y: 209
          width: 59
          height: 1
          color: "grey"
        }

        Slider {
          id: ch4
          anchors.horizontalCenter: parent.horizontalCenter
          y: 118
          width: 69
          height: 329
          to: 100
          orientation: Qt.Vertical
          value: strip4.faderValue
          enabled: !lockButton.checked
          onMoved: strip4.fadermoved(position)
          Material.accent: enabled ? defaultAccent : defaultDisabled
        }

        Button {
            id: mute4
            anchors.horizontalCenter: parent.horizontalCenter
            anchors.topMargin: 10
            anchors.top: parent.top
            width: 71
            height: 60
            text: qsTr("mute")
            scale: 1
            checkable: true
            checked: strip4.muted
            onToggled: strip4.muteclicked(checked)
        }

        Text {
            anchors.horizontalCenter: parent.horizontalCenter
            anchors.top: ch4.bottom
            anchors.topMargin: 10
            text: strip4.label
            font.pixelSize: 14
            horizontalAlignment: TextInput.AlignHCenter
        }

        ProgressBar {
          value: strip4.meterValue
          width: 61
          anchors.horizontalCenter: parent.horizontalCenter
          anchors.top: mute4.bottom
          anchors.topMargin: 15
        }
      }

      Item {
        id: chStrip5
        anchors.topMargin: 0
        anchors.top: parent.top
        anchors.left: chStrip4.right
        width: 114
        visible: configuration.channelStrips[5].enabled
        enabled: configuration.channelStrips[5].enabled

        Rectangle {
          anchors.horizontalCenter: parent.horizontalCenter
          y: 209
          width: 59
          height: 1
          color: "grey"
        }

        Slider {
          id: ch5
          anchors.horizontalCenter: parent.horizontalCenter
          y: 118
          width: 69
          height: 329
          to: 100
          orientation: Qt.Vertical
          value: strip5.faderValue
          enabled: !lockButton.checked
          onMoved: strip5.fadermoved(position)
          Material.accent: enabled ? defaultAccent : defaultDisabled
        }

        Button {
            id: mute5
            anchors.horizontalCenter: parent.horizontalCenter
            anchors.topMargin: 10
            anchors.top: parent.top
            width: 71
            height: 60
            text: qsTr("mute")
            scale: 1
            checkable: true
            checked: strip5.muted
            onToggled: strip5.muteclicked(checked)
        }

        Text {
            anchors.horizontalCenter: parent.horizontalCenter
            anchors.top: ch5.bottom
            anchors.topMargin: 10
            text: strip5.label
            font.pixelSize: 14
            horizontalAlignment: TextInput.AlignHCenter
        }

        ProgressBar {
          value: strip5.meterValue
          width: 61
          anchors.horizontalCenter: parent.horizontalCenter
          anchors.top: mute5.bottom
          anchors.topMargin: 15
        }
      }

      Item {
        id: buttonArea
        anchors.topMargin: 0
        anchors.top: parent.top
        anchors.left: chStrip5.right
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

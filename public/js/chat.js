/**
 * chat demo
 * Created by phachon@163.com
 */

var Chat = {

	/**
	 * ws url
	 */
	wsUrl : 'ws://' + window.location.host + '/ws',

	/**
	 * user_id
	 */
	userId: 0,

	/**
	 * socket
	 */
	socket: null,

	/**
	 * chat div element
	 */
	chatElement: '#chat_box',

	/**
	 * 初始化，连接服务端
	 */
	init: function () {
		var userId = Math.floor((Math.random()*1000) + 1);
		var roomId = Math.floor((Math.random()*1000) + 1);

		$("input[name='room_id']").val(roomId);
		$("input[name='user_id']").val(userId);

		$("input[name='room_id']").attr("disabled",true);
		$("input[name='user_id']").attr("disabled",true);
		this.socket = new WebSocket(this.wsUrl + "?room_id="+ roomId +"&user_id=" + userId);
		this.bindSocket();
	},

	bindSocket: function () {
		//socket 连接
		var socket = this.socket;
		if(socket == null) {
			this.init();
		}
		//连接成功
		socket.onopen = function (event) {
			onOpen(event);
		};
		//连接关闭
		socket.onclose = function (event) {
			onClose(event);
		};
		//返回消息
		socket.onmessage = function (event) {
			onMessage(event);
		};
		//错误
		socket.onerror = function (event) {
			onError(event);
		};

		function onOpen(event) {
			writeToScreen("websocket connect success!");
		}

		function onClose(event) {
			writeToScreen("websocket connect close!");
		}

		function onMessage(event) {
			writeToScreen(event.data);
		}

		function onError(event) {
			writeToScreen('error:'+ event.data);
		}

		function writeToScreen(messages) {
			console.log(messages);
			messages = eval("("+messages+")");
			var html = '<div class="alert alert-info alert-dismissible fade in" role="alert" style="padding: 10px 10px;margin-bottom: 10px;">';
			html += '<p>';
			html += '<span class="glyphicon glyphicon-volume-up" aria-hidden="true"></span>';
			html += '<strong> '+messages.user_id+'：</strong>';
			html += messages.message;
			html += '</p>';
			html += '</div>';
			$("#chat_box").append(html);
			$('#chat_box').scrollTop( $('#chat_box')[0].scrollHeight);
		}
	},

	send: function () {
		var socket = this.socket;
		if(socket == null) {
			this.init();
		}
		var id = $("input[name='user_id']").val();
		var message = $("textarea[name='message']").val();
		var data = '{"user_id": "'+ id +'", "type": 1, "message":"'+ message +'"}';
		// alert(data);
		socket.send(data);

		$("textarea[name='message']").val("");
	}

};

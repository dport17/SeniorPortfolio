<?php
	session_start();    
	if (checklogin_mysql($_POST["username"],$_POST["password"])) {
?>
	<h2> Welcome <?php echo $_POST['username']; ?> !</h2>
<?php		
	}else{
		echo "<script>alert('Invalid username/password');</script>";
		die();
	}
	function checklogin($username, $password) {
		$account = array("admin","1234");
		if (($username== $account[0]) and ($password == $account[1])) 
		  return TRUE;
		else return FALSE;
  	}
  	function checklogin_mysql($username, $password) {
  		$mysqli = new mysql('localhost','porterd3','***', 'secad_porterd3')
		$sql = "SELECT * FROM users WHERE username='" . $username . "' AND password = password('" . $password ."')";
		echo "DEBUG>sql=$sql"; 
		return TRUE;
  	}

?>

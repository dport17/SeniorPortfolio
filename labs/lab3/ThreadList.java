import java.util.ArrayList;
class ThreadList{
	private ArrayList<EchoServerThread> threadlist = new ArrayList<EchoServerThread>();
	public ThreadList(){	

	}
	public int getNumberofThreads(){
		return threadlist.size();
	//return the number of current threads
	}
	public void addThread(EchoServerThread newthread){
		threadlist.add(newthread);
	//add the newthread object to the threadlist	
	}
	public void removeThread(EchoServerThread thread){
		threadlist.remove(thread);
	//remove the given thread from the threadlist		
	}
	public void sendToAll(String message){
	//ask each thread in the threadlist to send the given message to its client		
	}
}
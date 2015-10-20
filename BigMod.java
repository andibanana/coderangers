import java.util.Scanner;
import java.math.BigInteger;

public class Main{

    public static void main(String[] args){
        Scanner in = new Scanner(System.in);
        while(in.hasNext()){
            BigInteger B = new BigInteger(in.next());
            BigInteger P = new BigInteger(in.next());
            BigInteger M = new BigInteger(in.next());
            System.out.println(B.modPow(P, M));
        }
    }

}
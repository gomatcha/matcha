package matcha;

public class Bridge {
    private static final Bridge instance = new Bridge();
    private Bridge() {
    }
    public static Bridge singleton() {
        return instance;
    }
}
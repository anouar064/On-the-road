// Name: Anouar Smaili

public class Point3D {

    private double x;
    private double y;
    private double z;

    // Contruct a 3D point from x, y, z coordinates
    public Point3D(double x, double y, double z) {
        this.x = x;
        this.y = y;
        this.z = z;
    }

    // Getters
    public double getX() {
        return this.x;
    }

    public double getY() {
        return this.y;
    }

    public double getZ() {
        return this.z;
    }
}
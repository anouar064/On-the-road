// Name: Anouar Smaili

public class Plane3D {

    private double a;
    private double b;
    private double c;
    private double d;

    // Contructor of a 3D Plane from 3 points
    public Plane3D(Point3D p1, Point3D p2, Point3D p3) {
        // Calculations in order to get the plane equation
        double a1 = p2.getX() - p1.getX();
        double b1 = p2.getY() - p1.getY();
        double c1 = p2.getZ() - p1.getZ();
        double a2 = p3.getX() - p1.getX();
        double b2 = p3.getY() - p1.getY();
        double c2 = p3.getZ() - p1.getZ();
        this.a = b1 * c2 - b2 * c1;
        this.b = a2 * c1 - a1 * c2;
        this.c = a1 * b2 - b1 * a2;
        this.d = (-a * p1.getX() - b * p1.getY() - c * p1.getZ());
    }

    // Contructor of a 3D Plane using the plane equation
    public Plane3D(double a, double b, double c, double d) {
        this.a = a;
        this.b = b;
        this.c = c;
        this.d = d;
    }

    // Returns the distance from a point to the plane
    public double getDistance(Point3D pt) {
        // Formula
        double m = Math.abs((a * pt.getX() + b * pt.getY() + c * pt.getZ() + d));
        double n = (double) Math.sqrt(a * a + b * b + c * c);
        return m / n;
    }

}

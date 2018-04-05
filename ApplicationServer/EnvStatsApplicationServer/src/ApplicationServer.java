import java.io.IOException;
import java.sql.*;
import java.util.UUID;

import com.rabbitmq.client.*;

public class ApplicationServer {
    private static java.sql.Connection connect = null;
    private static Statement statement = null;
    private static PreparedStatement preparedStatement = null;
    private static ResultSet resultSet = null;
    
    private static final String DATABASE = "db1";
    private static final String MEASUREMENT = "measurement";
    private static final String SENSOR = "sensor";
    private static final String SENSORTYPE = "sensortype";

    private static final String ip = "94.23.200.26:3306";
    private static final String user = "netcomp";
    private static final String pw = "envstat";
    
    private static void connectToDatabase() throws SQLException, ClassNotFoundException{
        Class.forName("com.mysql.jdbc.Driver");
        // Setup the connection with the DB
        connect = DriverManager
                .getConnection("jdbc:mysql://"+ip, user, pw);

        // Statements allow to issue SQL queries to the database
        statement = connect.createStatement();
    }

    public static void getColumnNames(ResultSet rs) throws SQLException {
        if (rs == null) {
            return;
        }
        ResultSetMetaData rsMetaData = rs.getMetaData();
        int numberOfColumns = rsMetaData.getColumnCount();

        // get the column names; column indexes start from 1
        for (int i = 1; i < numberOfColumns + 1; i++) {
// wtf

            String columnName = rsMetaData.getColumnName(i);
            // Get the name of the column's table name
            String tableName = rsMetaData.getTableName(i);
            System.out.println("column name=" + columnName + " table=" + tableName + "");
        }
    }
    
    private static void insertData(long mac, long sensorTimestamp, float latitude, float longitude, String sensorType, float data) throws SQLException {
        // insert SensorType

        String uuid = UUID.randomUUID().toString();
        long serverTimestamp = new Timestamp(System.currentTimeMillis()).getTime();

//        String query = "select id, name, radius from " + DATABASE+"."+ SENSORTYPE+"";
//        // create a statement
//        statement = connect.createStatement();
//        // execute query and return result as a ResultSet
//        resultSet = statement.executeQuery(query);
//
//        getColumnNames(resultSet);

        preparedStatement = connect
                .prepareStatement("INSERT INTO " + DATABASE+"."+ SENSORTYPE + " (name, radius) VALUES (?, ?);");
        preparedStatement.setString(1, sensorType);
        preparedStatement.setInt(2, 9);
        preparedStatement.executeUpdate();

        // insert Sensor
        preparedStatement = connect
                .prepareStatement("insert into " + DATABASE+"."+ SENSOR + " values (?, ?, ?, ?, ?)");
        preparedStatement.setString(1, uuid );
        preparedStatement.setLong(2, mac);
        preparedStatement.setInt(3, 0);
        preparedStatement.setFloat(4, latitude);
        preparedStatement.setFloat(5, longitude);
        preparedStatement.executeUpdate();

        // insert measurement
        preparedStatement = connect
                .prepareStatement("insert into " + DATABASE+"."+ MEASUREMENT + " values (?, ?, ?, ?)");
        preparedStatement.setString(1, UUID.randomUUID().toString() );
        preparedStatement.setString(2, uuid );
        preparedStatement.setFloat(3, data);
        //preparedStatement.setLong(5, sensorTimestamp);
        preparedStatement.setLong(4, serverTimestamp);
        preparedStatement.executeUpdate();
    }
    
    private static void close() {
        try {
            if (resultSet != null) {
                resultSet.close();
            }

            if (statement != null) {
                statement.close();
            }

            if (connect != null) {
                connect.close();
            }
        } catch (Exception e) {

        }
    }

    /**
     * The channel that the server listens on.
     */
    private static final String QUEUE_NAME = "sensor_data";

    private static boolean running = true;

    public static void main(String[] args) throws Exception {
        System.out.println("The server is running.");

        ConnectionFactory factory = new ConnectionFactory();
        factory.setHost("localhost");
        com.rabbitmq.client.Connection connection = factory.newConnection();
        Channel channel = connection.createChannel();

        channel.queueDeclare(QUEUE_NAME, false, false, false, null);

        System.out.println("Connecting to Database...");
        try {
            connectToDatabase();
            System.out.println(" [*] Waiting for messages. To exit press CTRL+C");

            Consumer consumer = new DefaultConsumer(channel) {
                @Override
                public void handleDelivery(String consumerTag, Envelope envelope, AMQP.BasicProperties properties, byte[] body)
                        throws IOException {
                    String message = new String(body, "UTF-8");
                    System.out.println(" [x] Received '" + message + "'");
                    System.out.println("Trying to insert data: " + message);
                    String[] values = message.split("|");
                    try {
                        insertData(Long.parseLong(values[0]), Long.parseLong(values[1]), Float.parseFloat(values[2]),
                                Float.parseFloat(values[3]), values[4], Float.parseFloat(values[5]));
                    } catch (SQLException e) {
                        e.printStackTrace();
                    }
                }
            };
            while(running)
                channel.basicConsume(QUEUE_NAME, true, consumer);

        } catch (Exception e){
            System.out.println("Failed to connect to Database:\n"+e);
             //do something
        } finally {
            close();
        }
    }
}
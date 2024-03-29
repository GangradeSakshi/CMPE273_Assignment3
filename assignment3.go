package main

import (
    "fmt"
    "github.com/julienschmidt/httprouter"
    "net/http"
    "encoding/json"
    "strings"
    "io/ioutil"
    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
    "strconv"
    "sort"
    "bytes"
)


type PriceEstimates struct {
    StartLatitude  float64
    StartLongitude float64
    EndLatitude    float64
    EndLongitude   float64
    Prices         []PriceEstimate `json:"prices"`
}


type PriceEstimate struct {
    ProductId       string  `json:"product_id"`
    CurrencyCode    string  `json:"currency_code"`
    DisplayName     string  `json:"display_name"`
    Estimate        string  `json:"estimate"`
    LowEstimate     int     `json:"low_estimate"`
    HighEstimate    int     `json:"high_estimate"`
    SurgeMultiplier float64 `json:"surge_multiplier"`
    Duration        int     `json:"duration"`
    Distance        float64 `json:"distance"`
}


func (pe *PriceEstimates) get(c *Client) error {
    priceEstimateParams := map[string]string{
        "start_latitude":  strconv.FormatFloat(pe.StartLatitude, 'f', 2, 32),
        "start_longitude": strconv.FormatFloat(pe.StartLongitude, 'f', 2, 32),
        "end_latitude":    strconv.FormatFloat(pe.EndLatitude, 'f', 2, 32),
        "end_longitude":   strconv.FormatFloat(pe.EndLongitude, 'f', 2, 32),
    }

    data := c.getRequest("estimates/price", priceEstimateParams)
    if e := json.Unmarshal(data, &pe); e != nil {
        return e
    }
    return nil
}

const (
  
    APIUrl string = "https://sandbox-api.uber.com/v1/%s%s"
)


type Getdetails interface {
    get(c *Client) error
}

// OAuth parameters
type RequestOptions struct {
    ServerToken    string
    ClientId       string
    ClientSecret   string
    AppName        string
    AuthorizeUrl   string
    AccessTokenUrl string
    AccessToken string
    BaseUrl        string
}

// Client contains the required OAuth tokens and urls and manages
// the connection to the API. All requests are made via this type
type Client struct {
    Options *RequestOptions
}

// Create returns a new API client
func Create(options *RequestOptions) *Client {
    return &Client{options}
}

// Get formulates an HTTP GET request based on the Uber endpoint type
func (c *Client) Get(Getdetails Getdetails) error {
    if e := Getdetails.get(c); e != nil {
        return e
    }

    return nil
}


func (c *Client) getRequest(endpoint string, params map[string]string) []byte {
    urlParams := "?"
    params["server_token"] = c.Options.ServerToken
    for k, v := range params {
        if len(urlParams) > 1 {
            urlParams += "&"
        }
        urlParams += fmt.Sprintf("%s=%s", k, v)
    }

    url := fmt.Sprintf(APIUrl, endpoint, urlParams)

    res, err := http.Get(url)
    if err != nil {
      
    }

    data, err := ioutil.ReadAll(res.Body)
    res.Body.Close()

    return data
}

type Products struct {
    Latitude  float64
    Longitude float64
    Products  []Product `json:"products"`
}


type Product struct {
    ProductId   string `json:"product_id"`
    Description string `json:"description"`
    DisplayName string `json:"display_name"`
    Capacity    int    `json:"capacity"`
    Image       string `json:"image"`
}


func (product *Products) get(c *Client) error {
    productParams := map[string]string{
        "latitude":  strconv.FormatFloat(product.Latitude, 'f', 2, 32),
        "longitude": strconv.FormatFloat(product.Longitude, 'f', 2, 32),
    }

    data := c.getRequest("products", productParams)
    if e := json.Unmarshal(data, &product); e != nil {
        return e
    }
    return nil
}



type reqObj struct{
Id int
Name string `json:"Name"`
Address string `json:"Address"`
City string `json:"City"`
State string `json:"State"`
Zip string `json:"Zip"`
Coordinates struct{
    Lat float64
    Lng float64
}
}

var id int;
var tripId int;


type Responz struct {
    Results []struct {
        AddressComponents []struct {
            LongName  string   `json:"long_name"`
            ShortName string   `json:"short_name"`
            Types     []string `json:"types"`
        } `json:"address_components"`
        FormattedAddress string `json:"formatted_address"`
        Geometry         struct {
            Location struct {
                Lat float64 `json:"lat"`
                Lng float64 `json:"lng"`
            } `json:"location"`
            LocationType string `json:"location_type"`
            Viewport     struct {
                Northeast struct {
                    Lat float64 `json:"lat"`
                    Lng float64 `json:"lng"`
                } `json:"northeast"`
                Southwest struct {
                    Lat float64 `json:"lat"`
                    Lng float64 `json:"lng"`
                } `json:"southwest"`
            } `json:"viewport"`
        } `json:"geometry"`
        PartialMatch bool     `json:"partial_match"`
        PlaceID      string   `json:"place_id"`
        Types        []string `json:"types"`
    } `json:"results"`
    Status string `json:"status"`
}

type TripResponse struct {
    BestRouteLocationIds   []string `json:"best_route_location_ids"`
    ID                     string   `json:"id"`
    StartingFromLocationID string   `json:"starting_from_location_id"`
    Status                 string   `json:"status"`
    TotalDistance          float64  `json:"total_distance"`
    TotalUberCosts         int      `json:"total_uber_costs"`
    TotalUberDuration      int      `json:"total_uber_duration"`
}

type RideRequest struct {
    EndLatitude    string `json:"end_latitude"`
    EndLongitude   string `json:"end_longitude"`
    ProductID      string `json:"product_id"`
    StartLatitude  string `json:"start_latitude"`
    StartLongitude string `json:"start_longitude"`
}

type OnGoingTrip struct {
    BestRouteLocationIds      []string `json:"best_route_location_ids"`
    ID                        string   `json:"id"`
    NextDestinationLocationID string   `json:"next_destination_location_id"`
    StartingFromLocationID    string   `json:"starting_from_location_id"`
    Status                    string   `json:"status"`
    TotalDistance             float64  `json:"total_distance"`
    TotalUberCosts            int      `json:"total_uber_costs"`
    TotalUberDuration         int      `json:"total_uber_duration"`
    UberWaitTimeEta           int      `json:"uber_wait_time_eta"`
}

type ReqResponse struct {
    Driver          interface{} `json:"driver"`
    Eta             int         `json:"eta"`
    Location        interface{} `json:"location"`
    RequestID       string      `json:"request_id"`
    Status          string      `json:"status"`
    SurgeMultiplier int         `json:"surge_multiplier"`
    Vehicle         interface{} `json:"vehicle"`
}




type requestObject struct{
Greeting string
}

func createlocation(rw http.ResponseWriter, req *http.Request, p httprouter.Params) {
    id=id+1;


    decoder := json.NewDecoder(req.Body)
    var t reqObj 
    t.Id = id; 
    err := decoder.Decode(&t)
    if err != nil {
        fmt.Println("Error in decoding values")
    }

    st:=strings.Join(strings.Split(t.Address," "),"+");
    fmt.Println(st);
    constr := []string {strings.Join(strings.Split(t.Address," "),"+"),strings.Join(strings.Split(t.City," "),"+"),t.State}
    lstringplus := strings.Join(constr,"+")
    locstr := []string{"http://maps.google.com/maps/api/geocode/json?address=",lstringplus}

    resp, err := http.Get(strings.Join(locstr,""))
  
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
       fmt.Println("Error: Wrong address");
     }
     var data Responz
    err = json.Unmarshal(body, &data)
    fmt.Println(data.Status)
   
    t.Coordinates.Lat=data.Results[0].Geometry.Location.Lat;
    t.Coordinates.Lng=data.Results[0].Geometry.Location.Lng;




 conn, err := mgo.Dial("mongodb://admin:admin@ds045064.mongolab.com:45064/tripadvisor")


    if err != nil {
        panic(err)
    }
    defer conn.Close();

conn.SetMode(mgo.Monotonic,true);
c:=conn.DB("tripadvisor").C("location");
err = c.Insert(t);


    js,err := json.Marshal(t)
    if err != nil{
	   fmt.Println("Error")
	   return
	}
    rw.Header().Set("Content-Type","application/json")
    rw.Write(js)
}

func getlocation(rw http.ResponseWriter, req *http.Request, p httprouter.Params){
fmt.Println(p.ByName("locid"));
id ,err1:= strconv.Atoi(p.ByName("locid"))
if err1 != nil {
        panic(err1)
    }
 conn, err := mgo.Dial("mongodb://admin:admin@ds045064.mongolab.com:45064/tripadvisor")

    
    if err != nil {
        panic(err)
    }
    defer conn.Close();

conn.SetMode(mgo.Monotonic,true);
c:=conn.DB("tripadvisor").C("location");
result:=reqObj{}
err = c.Find(bson.M{"id":id}).One(&result)
if err != nil {
                fmt.Println(err)
        }

        js,err := json.Marshal(result)
    if err != nil{
       fmt.Println("Error")
       return
    }
    rw.Header().Set("Content-Type","application/json")
    rw.Write(js)
}

type modReqObj struct{
    Address string `json:"address"`
    City string `json:"city"`
    State string `json:"state"`
    Zip string `json:"zip"`
}

func updateLocation(rw http.ResponseWriter, req *http.Request, p httprouter.Params){
    
 id ,err1:= strconv.Atoi(p.ByName("locid"))
 //fmt.Println(id);
 if err1 != nil {
         panic(err1)
     }
  conn, err := mgo.Dial("mongodb://admin:admin@ds045064.mongolab.com:45064/tripadvisor")

//     // Check if connection error, is mongo running?
     if err != nil {
         panic(err)
     }
     defer conn.Close();

conn.SetMode(mgo.Monotonic,true);
 c:=conn.DB("tripadvisor").C("location");


     decoder := json.NewDecoder(req.Body)
     var t modReqObj  
     err = decoder.Decode(&t)
     if err != nil {
         fmt.Println("Error")
     }


     colQuerier := bson.M{"id": id}
     change := bson.M{"$set": bson.M{"address": t.Address, "city":t.City,"state":t.State,"zip":t.Zip}}
     err = c.Update(colQuerier, change)
     if err != nil {
         panic(err)
     }

}

func deletelocation(rw http.ResponseWriter, req *http.Request, p httprouter.Params){
     id ,err1:= strconv.Atoi(p.ByName("locid"))
 //fmt.Println(id);
 if err1 != nil {
         panic(err1)
     }
  conn, err := mgo.Dial("mongodb://admin:admin@ds045064.mongolab.com:45064/tripadvisor")
  conn.SetMode(mgo.Monotonic,true);
c:=conn.DB("tripadvisor").C("location");

//     // Check if connection error, is mongo running?
     if err != nil {
         panic(err)
     }
     defer conn.Close();
     err=c.Remove(bson.M{"id":id})
     if err != nil { fmt.Printf("Could not find kitten %s to delete", id)}
    rw.WriteHeader(http.StatusNoContent)
}

type userUber struct {
    LocationIds            []string `json:"location_ids"`
    StartingFromLocationID string   `json:"starting_from_location_id"`
}

func plantrip(rw http.ResponseWriter, req *http.Request, p httprouter.Params){

    decoder := json.NewDecoder(req.Body)
    var uUD userUber 
    err := decoder.Decode(&uUD)
    if err != nil {
        fmt.Println("Error")
    }

        fmt.Println(uUD.StartingFromLocationID);



    var options RequestOptions;
    options.ServerToken= "S37TkXJu1TBNbDea22MxgIAjoM1C__fJ3r6vbQ-5";
    options.ClientId= "5-BNiHDpt1CZvQoWd2G2vV2GSvSnIu2j";
    options.ClientSecret= "P5qyGJI-sJw5m-s2kFHljzg59kccexZ8qkbaL44P";
    options.AppName= "CMPE273-A3";
    options.BaseUrl= "https://sandbox-api.uber.com/v1/";
    

    client :=Create(&options); 

        sid ,err1:= strconv.Atoi(uUD.StartingFromLocationID)

 if err1 != nil {
         panic(err1)
     }

    conn, err := mgo.Dial("mongodb://admin:admin@ds045064.mongolab.com:45064/tripadvisor");

    if err != nil {
        panic(err)
    }
    defer conn.Close();

    conn.SetMode(mgo.Monotonic,true);
    c:=conn.DB("tripadvisor").C("location");
    result:=reqObj{}
    err = c.Find(bson.M{"id":sid}).One(&result)
    if err != nil {
                fmt.Println(err)
        }

    index:=0;
    totalPrice := 0;
    totalDistance :=0.0;
    totalDuration :=0;
    bestroute:=make([]float64,len(uUD.LocationIds));
    m := make(map[float64]string)

    for _,ids := range uUD.LocationIds{
    
        lid,err1:= strconv.Atoi(ids)
            //fmt.Println(id);
        if err1 != nil {
            panic(err1)
        }
        

        resultLID:=reqObj{}
        err = c.Find(bson.M{"id":lid}).One(&resultLID)
        if err != nil {
             fmt.Println(err)
        }
        pe := &PriceEstimates{}
        pe.StartLatitude = result.Coordinates.Lat;
        pe.StartLongitude = result.Coordinates.Lng;
        pe.EndLatitude = resultLID.Coordinates.Lat;
        pe.EndLongitude = resultLID.Coordinates.Lng;

        if e := client.Get(pe); e != nil {
            fmt.Println(e);
        }
        totalDistance=totalDistance+pe.Prices[0].Distance;
        totalDuration=totalDuration+pe.Prices[0].Duration;
        totalPrice=totalPrice+pe.Prices[0].LowEstimate;
        bestroute[index]=pe.Prices[0].Distance;
        m[pe.Prices[0].Distance]=ids;
        index=index+1;
    }
   
    sort.Float64s(bestroute);
    

    var tripres TripResponse;

    tripId=tripId+1;

     tripres.ID=strconv.Itoa(tripId);
     tripres.TotalDistance=totalDistance;
     tripres.TotalUberCosts=totalPrice;
     tripres.TotalUberDuration=totalDuration;
     tripres.Status="Planning";
     tripres.StartingFromLocationID=strconv.Itoa(sid);
     tripres.BestRouteLocationIds=make([]string,len(uUD.LocationIds));
     index=0;
     for _, ind := range bestroute{
        tripres.BestRouteLocationIds[index]=m[ind];
        index=index+1;
     }
     fmt.Println(tripres.BestRouteLocationIds[1]);

    

    c1:=conn.DB("tripadvisor").C("usertrips");
    err = c1.Insert(tripres);

     
        js,err := json.Marshal(tripres)
    if err != nil{
       fmt.Println("Error")
       return
    }
    rw.Header().Set("Content-Type","application/json")
    rw.Write(js)

    }


func gettrip(rw http.ResponseWriter, req *http.Request, p httprouter.Params){

    conn, err := mgo.Dial("mongodb://admin:admin@ds045064.mongolab.com:45064/tripadvisor")

    
    if err != nil {
        panic(err)
    }
    defer conn.Close();

    conn.SetMode(mgo.Monotonic,true);
    c:=conn.DB("tripadvisor").C("usertrips");
    result:=TripResponse{}
    err = c.Find(bson.M{"id":p.ByName("tripid")}).One(&result)
    if err != nil {
        fmt.Println(err)
    }

    
    js,err := json.Marshal(result)
    if err != nil{
       fmt.Println("Error in marshalling")
       return
    }
    rw.Header().Set("Content-Type","application/json")
    rw.Write(js)
}


var currentPos int;
var ogtID int;



func requestfortrip(rw http.ResponseWriter, req *http.Request, p httprouter.Params){
    
    kid ,err1:= strconv.Atoi(p.ByName("tripid"))
    var siD int;
    //fmt.Println(id);
    if err1 != nil {
         panic(err1)
     }
    var ogt OnGoingTrip;
    result1:=reqObj{}
    result2:=reqObj{}
    conn, err := mgo.Dial("mongodb://admin:admin@ds045064.mongolab.com:45064/tripadvisor")

   
    if err != nil {
        panic(err)
    }
    defer conn.Close();

    conn.SetMode(mgo.Monotonic,true);
    c:=conn.DB("tripadvisor").C("usertrips");
    result:=TripResponse{}

    err = c.Find(bson.M{"id":strconv.Itoa(kid)}).One(&result)
    if err != nil {
        fmt.Println(err)
    }else{

    var iD int;

    c1:=conn.DB("tripadvisor").C("location");
    if currentPos==0{
        iD, err = strconv.Atoi(result.StartingFromLocationID)
        siD=iD;
        if err != nil {
       
            fmt.Println(err)
        }
    }else
    {
        iD, err = strconv.Atoi(result.BestRouteLocationIds[currentPos-1])
        siD=iD;
        if err != nil {
       
            fmt.Println(err)
        }
    }

    err = c1.Find(bson.M{"id":iD}).One(&result1)
    if err != nil {
                fmt.Println(err)
        }
    iD, err = strconv.Atoi(result.BestRouteLocationIds[currentPos])
    if err != nil {
        // handle error
        fmt.Println(err)
    }
    err = c1.Find(bson.M{"id":iD}).One(&result2)
    if err != nil {
                fmt.Println(err)
        }


        fmt.Println(result2.Coordinates.Lat);
    }

    ogt.ID=strconv.Itoa(ogtID);
    ogt.BestRouteLocationIds=result.BestRouteLocationIds;
    ogt.StartingFromLocationID=strconv.Itoa(siD);
    ogt.NextDestinationLocationID=result.BestRouteLocationIds[currentPos];
    ogt.TotalDistance=result.TotalDistance;
    ogt.TotalUberCosts=result.TotalUberCosts;
    ogt.TotalUberDuration=result.TotalUberDuration;
    ogt.Status="requesting";
    


    var options RequestOptions;
    options.ServerToken= "3yWU14IA8z2dhhfzk2LWtMe0eKcbessaQVDPrQ9l";
    options.ClientId= "_fxzbWklwtQn-xaC3vcHtO5wdFr_2AsB";
    options.ClientSecret= "2a1YHPslNpqLovztEa05ODKx0KndzK_l19p22x5g";
    options.AppName= "MyTripAdvisor";
    options.BaseUrl= "https://sandbox-api.uber.com/v1/";
    client :=Create(&options);

    pl:=Products{};
    pl.Latitude=result1.Coordinates.Lat;
    pl.Longitude=result1.Coordinates.Lng;
    if e := pl.get(client); e != nil {
         fmt.Println(e)
    }
    var prodid string;
    i:=0
    for _, product := range pl.Products {
         if(i == 0){
             prodid = product.ProductId
        }
    }

    var rr RideRequest;

    rr.StartLatitude=strconv.FormatFloat(result1.Coordinates.Lat, 'f', 6, 64);
    rr.StartLongitude=strconv.FormatFloat(result1.Coordinates.Lng, 'f', 6, 64);
    rr.EndLatitude=strconv.FormatFloat(result2.Coordinates.Lat, 'f', 6, 64);
    rr.EndLongitude=strconv.FormatFloat(result2.Coordinates.Lng, 'f', 6, 64);
    rr.ProductID=prodid;
    buf, _ := json.Marshal(rr)
    body := bytes.NewBuffer(buf)
    url := fmt.Sprintf(APIUrl, "requests?","access_token=eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJzY29wZXMiOlsicmVxdWVzdCJdLCJzdWIiOiI5MjY3N2VlZC04MWIxLTQwYzctYmY4ZC1kOTMxMDQxMmIxZTkiLCJpc3MiOiJ1YmVyLXVzMSIsImp0aSI6IjY4MGViNjFmLWFjNjgtNDI3YS04YWM1LTgxMjcyNzFlZGQyNSIsImV4cCI6MTQ1MDc1Mzk3NywiaWF0IjoxNDQ4MTYxOTc3LCJ1YWN0IjoiZFppS2FaMGhXcGl6R25DOVAwU2prS3U5UlZLNWFrIiwibmJmIjoxNDQ4MTYxODg3LCJhdWQiOiJfZnh6YldrbHd0UW4teGFDM3ZjSHRPNXdkRnJfMkFzQiJ9.gmkvmQG081LF6wXWgM-5I1mrHYGgRo_3DUgxmDpO_xo3AlzyvQm5Vl5Ro2BAbxG-9HGmoXm5D-tbt0H6iT2x_tLw2zDkfU2wJ6-8GB-Vjgx4rEz974G6ktfB8I0I7pq_4ATCbds7J3ayLOH0vqKTExodtWY8N_onUMCuf5GsvwL9UcDKegpG8Oq0gZLcXYeMMywmuC6ZJu6s0MgGObqLAHJRZTbm5e2hDaOY_5pBHLG_PkIML7WOQ3UAWIUuWLBVeQ2GjCYY4NIipd9xyozXniK7d1GDJWKHhU0Fe1wn6-QPvoKt87wGX70HK-Ht36zer7PMlUqzLtpSF4LEfYf3eg")
    res, err := http.Post(url,"application/json",body)
    if err != nil {
        fmt.Println(err)
    }
    data, err := ioutil.ReadAll(res.Body)
    var rRes ReqResponse;
    err = json.Unmarshal(data, &rRes)
    ogt.UberWaitTimeEta=rRes.Eta;

    js,err := json.Marshal(ogt)
    if err != nil{
       fmt.Println("Error")
       return
    }
    ogtID=ogtID+1;
    currentPos=currentPos+1;
    rw.Header().Set("Content-Type","application/json")
    rw.Write(js)

}


func main() {
    mux := httprouter.New()
    
    id=0;
    tripId=0;
    currentPos=0;
    ogtID=0;
    mux.POST("/locations",createlocation)
    mux.POST("/trips",plantrip)
    mux.GET("/locations/:locid",getlocation)
    mux.GET("/trips/:tripid",gettrip)
    mux.PUT("/locations/:locid",updateLocation)
    mux.GET("/trips/:tripid/request",requestfortrip)
    mux.DELETE("/locations/:locid",deletelocation)
    server := http.Server{
            Addr:        "0.0.0.0:8080",
            Handler: mux,
    }   

    server.ListenAndServe()
}
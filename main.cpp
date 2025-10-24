#include <pistache/endpoint.h>
#include <pistache/router.h>
#include <nlohmann/json.hpp>
#include <vector>
#include <string>
#include <algorithm>
#include <iostream>
#include <sstream>
#include <unordered_map>

using namespace Pistache;
using json = nlohmann::json;

struct Record {
    int id;
    std::string first_name;
    std::string middle_name;
    std::string last_name;
    std::string street;
    std::string city;
    std::string state;
    std::string zip;
    std::string phone;
    std::string email;

    json to_json() const {
        return json{
            {"id", id},
            {"first_name", first_name},
            {"middle_name", middle_name},
            {"last_name", last_name},
            {"street", street},
            {"city", city},
            {"state", state},
            {"zip", zip},
            {"phone", phone},
            {"email", email}
        };
    }
};

class ApiHandler {
public:
    explicit ApiHandler(std::vector<Record>& records, int& next_id)
        : records_(records), next_id_(next_id) {}

    void create(const Rest::Request& request, Http::ResponseWriter response) {
        std::cout << "[POST /records] Creating record" << std::endl;
        try {
            auto body = json::parse(request.body());
            Record r;
            r.id = next_id_++;
            r.first_name = body.value("first_name", "");
            r.middle_name = body.value("middle_name", "");
            r.last_name = body.value("last_name", "");
            r.street = body.value("street", "");
            r.city = body.value("city", "");
            r.state = body.value("state", "");
            r.zip = body.value("zip", "");
            r.phone = body.value("phone", "");
            r.email = body.value("email", "");

            records_.push_back(r);
            std::cout << "[POST /records] Created record ID: " << r.id << std::endl;
            response.send(Http::Code::Created, r.to_json().dump());
        } catch (const std::exception& e) {
            std::cout << "[POST /records] ERROR: Invalid JSON - " << e.what() << std::endl;
            response.send(Http::Code::Bad_Request, "Invalid JSON");
        }
    }

    void read(const Rest::Request& request, Http::ResponseWriter response) {
        int id = request.param(":id").as<int>();
        std::cout << "[GET /records/" << id << "] Reading record" << std::endl;
        
        auto it = std::find_if(records_.begin(), records_.end(),
                               [id](const Record& r) { return r.id == id; });
        if (it == records_.end()) {
            std::cout << "[GET /records/" << id << "] ERROR: Record not found" << std::endl;
            response.send(Http::Code::Not_Found, "Record not found");
            return;
        }
        std::cout << "[GET /records/" << id << "] Found record" << std::endl;
        response.send(Http::Code::Ok, it->to_json().dump());
    }

    void update(const Rest::Request& request, Http::ResponseWriter response) {
        int id = request.param(":id").as<int>();
        std::cout << "[PUT /records/" << id << "] Updating record" << std::endl;
        
        try {
            auto body = json::parse(request.body());
            auto it = std::find_if(records_.begin(), records_.end(),
                                   [id](const Record& r) { return r.id == id; });
            if (it == records_.end()) {
                std::cout << "[PUT /records/" << id << "] ERROR: Record not found" << std::endl;
                response.send(Http::Code::Not_Found, "Record not found");
                return;
            }

            // FIX: Use .value() instead of .contains() + .get()
            it->first_name = body.value("first_name", it->first_name);
            it->middle_name = body.value("middle_name", it->middle_name);
            it->last_name = body.value("last_name", it->last_name);
            it->street = body.value("street", it->street);
            it->city = body.value("city", it->city);
            it->state = body.value("state", it->state);
            it->zip = body.value("zip", it->zip);
            it->phone = body.value("phone", it->phone);
            it->email = body.value("email", it->email);

            std::cout << "[PUT /records/" << id << "] Updated record" << std::endl;
            response.send(Http::Code::Ok, it->to_json().dump());
        } catch (const std::exception& e) {
            std::cout << "[PUT /records/" << id << "] ERROR: Invalid JSON - " << e.what() << std::endl;
            response.send(Http::Code::Bad_Request, "Invalid JSON");
        }
    }

    void del(const Rest::Request& request, Http::ResponseWriter response) {
        int id = request.param(":id").as<int>();
        std::cout << "[DELETE /records/" << id << "] Deleting record" << std::endl;
        
        auto it = std::find_if(records_.begin(), records_.end(),
                               [id](const Record& r) { return r.id == id; });
        if (it == records_.end()) {
            std::cout << "[DELETE /records/" << id << "] ERROR: Record not found" << std::endl;
            response.send(Http::Code::Not_Found, "Record not found");
            return;
        }
        records_.erase(it);
        std::cout << "[DELETE /records/" << id << "] Deleted record" << std::endl;
        response.send(Http::Code::No_Content, "");
    }

    void reset(const Rest::Request& request, Http::ResponseWriter response) {
        std::cout << "[DELETE /reset] Resetting database" << std::endl;
        records_.clear();
        next_id_ = 1;
        std::cout << "[DELETE /reset] Database cleared" << std::endl;
        response.send(Http::Code::No_Content, "");
    }


void query(const Rest::Request& request, Http::ResponseWriter response) {
    std::cout << "[GET /records] Flexible query started" << std::endl;
    json results = json::array();

    for (const auto& r : records_) {
        bool match = true;

        if (request.query().has("id")) {
            if (std::to_string(r.id) != request.query().get("id").value())
                match = false;
        }

        if (request.query().has("first_name")) {
            if (r.first_name != request.query().get("first_name").value())
                match = false;
        }

        if (request.query().has("middle_name")) {
            if (r.middle_name != request.query().get("middle_name").value())
                match = false;
        }

        if (request.query().has("last_name")) {
            if (r.last_name != request.query().get("last_name").value())
                match = false;
        }

        if (request.query().has("street")) {
            if (r.street != request.query().get("street").value())
                match = false;
        }

        if (request.query().has("city")) {
            if (r.city != request.query().get("city").value())
                match = false;
        }

        if (request.query().has("state")) {
            if (r.state != request.query().get("state").value())
                match = false;
        }

        if (request.query().has("zip")) {
            if (r.zip != request.query().get("zip").value())
                match = false;
        }

        if (request.query().has("phone")) {
            std::string queryPhone = request.query().get("phone").value();
            if (!(r.phone == queryPhone ||
                  (queryPhone.length() == 3 && r.phone.substr(0, 3) == queryPhone)))
                match = false;
        }

        if (request.query().has("email")) {
            if (r.email != request.query().get("email").value())
                match = false;
        }

        if (match)
            results.push_back(r.to_json());
    }

    std::cout << "[GET /records] Found " << results.size() << " matching records" << std::endl;
    response.send(Http::Code::Ok, results.dump());
}



private:
    std::vector<Record>& records_;
    int& next_id_;
};

int main() {
    std::vector<Record> records;
    int next_id = 1;

    std::cout << "Starting API server on http://localhost:8080" << std::endl;

    // Initialize Pistache server
    Http::Endpoint server(Address(Ipv4::any(), Port(8080)));
    auto opts = Http::Endpoint::options().threads(4);
    server.init(opts);

    // Set up router
    Rest::Router router;
    ApiHandler handler(records, next_id);

    // Define routes
    Rest::Routes::Post(router, "/records", Rest::Routes::bind(&ApiHandler::create, &handler));
    Rest::Routes::Get(router, "/records/:id", Rest::Routes::bind(&ApiHandler::read, &handler));
    Rest::Routes::Put(router, "/records/:id", Rest::Routes::bind(&ApiHandler::update, &handler));
    Rest::Routes::Delete(router, "/records/:id", Rest::Routes::bind(&ApiHandler::del, &handler));
    Rest::Routes::Delete(router, "/reset", Rest::Routes::bind(&ApiHandler::reset, &handler));
    Rest::Routes::Get(router, "/records", Rest::Routes::bind(&ApiHandler::query, &handler));

    // Start server
    server.setHandler(router.handler());
    server.serve();

    return 0;
}


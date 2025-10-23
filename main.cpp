#include <pistache/endpoint.h>
#include <pistache/router.h>
#include <nlohmann/json.hpp>
#include <vector>
#include <string>
#include <algorithm>

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
            response.send(Http::Code::Created, r.to_json().dump());
        } catch (const std::exception& e) {
            response.send(Http::Code::Bad_Request, "Invalid JSON");
        }
    }

    void read(const Rest::Request& request, Http::ResponseWriter response) {
        int id = request.param(":id").as<int>();
        auto it = std::find_if(records_.begin(), records_.end(),
                               [id](const Record& r) { return r.id == id; });
        if (it == records_.end()) {
            response.send(Http::Code::Not_Found, "Record not found");
            return;
        }
        response.send(Http::Code::Ok, it->to_json().dump());
    }

    void update(const Rest::Request& request, Http::ResponseWriter response) {
        int id = request.param(":id").as<int>();
        try {
            auto body = json::parse(request.body());
            auto it = std::find_if(records_.begin(), records_.end(),
                                   [id](const Record& r) { return r.id == id; });
            if (it == records_.end()) {
                response.send(Http::Code::Not_Found, "Record not found");
                return;
            }

            if (body.contains("first_name")) it->first_name = body["first_name"].get<std::string>();
            if (body.contains("middle_name")) it->middle_name = body["middle_name"].get<std::string>();
            if (body.contains("last_name")) it->last_name = body["last_name"].get<std::string>();
            if (body.contains("street")) it->street = body["street"].get<std::string>();
            if (body.contains("city")) it->city = body["city"].get<std::string>();
            if (body.contains("state")) it->state = body["state"].get<std::string>();
            if (body.contains("zip")) it->zip = body["zip"].get<std::string>();
            if (body.contains("phone")) it->phone = body["phone"].get<std::string>();
            if (body.contains("email")) it->email = body["email"].get<std::string>();

            response.send(Http::Code::Ok, it->to_json().dump());
        } catch (const std::exception& e) {
            response.send(Http::Code::Bad_Request, "Invalid JSON");
        }
    }

    void del(const Rest::Request& request, Http::ResponseWriter response) {
        int id = request.param(":id").as<int>();
        auto it = std::find_if(records_.begin(), records_.end(),
                               [id](const Record& r) { return r.id == id; });
        if (it == records_.end()) {
            response.send(Http::Code::Not_Found, "Record not found");
            return;
        }
        records_.erase(it);
        response.send(Http::Code::No_Content, "");
    }

    void reset(const Rest::Request& request, Http::ResponseWriter response) {
        records_.clear();
        next_id_ = 1;
        response.send(Http::Code::No_Content, "");
    }
  
    void query(const Rest::Request& request, Http::ResponseWriter response) {
        std::string first_name = request.query().get("first_name").value_or("");
        std::string middle_name = request.query().get("middle_name").value_or("");
        std::string last_name = request.query().get("last_name").value_or("");
        std::string street = request.query().get("street").value_or("");
        std::string city = request.query().get("city").value_or("");
        std::string state = request.query().get("state").value_or("");
        std::string zip = request.query().get("zip").value_or("");
        std::string phone = request.query().get("phone").value_or("");
        std::string email = request.query().get("email").value_or("");

        json results = json::array();
        for (const auto& r : records_) {
            bool match = true;

            if (!first_name.empty() && r.first_name != first_name) match = false;
            if (!middle_name.empty() && r.middle_name != middle_name) match = false;
            if (!last_name.empty() && r.last_name != last_name) match = false;
            if (!street.empty() && r.street != street) match = false;
            if (!city.empty() && r.city != city) match = false;
            if (!state.empty() && r.state != state) match = false;
            if (!zip.empty() && r.zip != zip) match = false;
            if (!email.empty() && r.email != email) match = false;

            if (!phone.empty()) {
                if (r.phone == phone) {
                    // Full match
                } else if (phone.length() == 3 && r.phone.length() >= 3 && r.phone.substr(0, 3) == phone) {
                    // Area code match
                } else {
                    match = false;
                }
            }

            if (match) {
                results.push_back(r.to_json());
            }
        }

        response.send(Http::Code::Ok, results.dump());
    }

private:
    std::vector<Record>& records_;
    int& next_id_;
};

int main() {
    std::vector<Record> records;
    int next_id = 1;

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


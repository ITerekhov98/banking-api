package cbr

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

const keyRateSOAPURL = "https://www.cbr.ru/DailyInfoWebServ/DailyInfo.asmx"

func buildKeyRateSOAPBody(targetDate time.Time) []byte {
	toDate := targetDate
	fromDate := targetDate.AddDate(0, 0, -1)
	layout := "2006-01-02T15:04:05"
	return fmt.Appendf(nil, `
<soap:Envelope xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
               xmlns:xsd="http://www.w3.org/2001/XMLSchema"
               xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <KeyRate xmlns="http://web.cbr.ru/">
      <fromDate>%s</fromDate>
      <ToDate>%s</ToDate>
    </KeyRate>
  </soap:Body>
</soap:Envelope>`, fromDate.Format(layout), toDate.Format(layout))
}

func GetKeyRate() (float64, error) {
	req, err := http.NewRequest("POST", keyRateSOAPURL, bytes.NewReader(buildKeyRateSOAPBody(time.Now())))
	if err != nil {
		return 0, err
	}

	req.Header.Set("Content-Type", "text/xml; charset=utf-8")
	req.Header.Set("SOAPAction", `"http://web.cbr.ru/KeyRate"`)

	client := http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return 0, fmt.Errorf("CBR request failed: %s\n%s", resp.Status, string(body))
	}

	return parseSOAPResponse(resp.Body)
}

func parseSOAPResponse(body io.Reader) (float64, error) {
	type KR struct {
		Date string `xml:"DT"`
		Rate string `xml:"Rate"`
	}

	type KeyRate struct {
		Items []KR `xml:"KR"`
	}

	type Diffgram struct {
		KeyRate KeyRate `xml:"KeyRate"`
	}

	type KeyRateResult struct {
		Diffgram Diffgram `xml:"diffgram"`
	}

	type Envelope struct {
		Body struct {
			Response struct {
				Result KeyRateResult `xml:"KeyRateResult"`
			} `xml:"KeyRateResponse"`
		} `xml:"Body"`
	}

	var env Envelope
	if err := xml.NewDecoder(body).Decode(&env); err != nil {
		return 0, fmt.Errorf("failed to decode XML: %w", err)
	}

	items := env.Body.Response.Result.Diffgram.KeyRate.Items
	if len(items) == 0 {
		return 0, fmt.Errorf("no key rate entries found")
	}

	// Возьмём последнюю ставку
	rawRate := items[len(items)-1].Rate
	return strconv.ParseFloat(rawRate, 64)
}

package services

import (
    "math/rand"
    "sort"
    "time"
    "tribo_ofertas_backend/models"
    "gorm.io/gorm"
)

func UpdateOffers(db *gorm.DB) error {
    var offers []models.Offer
    if err := db.Find(&offers).Error; err != nil {
        return err
    }

    for i := range offers {
        offer := &offers[i]

        // Verifica se a oferta está habilitada
        verificarValores(offer)
        if offer.Enable {
            // Calcula as variáveis intermediárias
            calcularVariaveisIntermediarias(offer)

            // Atualiza o valor de YNorm
            offer.YNorm = calcularY(1.0, 2.0, -1.0, 0.5, offer.R1, offer.R2, offer.R3, offer.R4)
        }
    }

    // Salva as ofertas atualizadas no banco de dados
    return db.Save(&offers).Error
}

// Função para transformar a diferença entre duas datas em um número decimal
func transfCountDateTime(duration time.Duration) float64 {
    return duration.Hours() / 24 // Convertendo para dias
}

// Função para verificar os valores obrigatórios
func verificarValores(offer *models.Offer) {
    if offer.Price <= 0 || offer.PromoUnits == 0 || offer.StartDate.IsZero() || offer.EndDate.IsZero() {
        offer.Enable = false
    } else {
        offer.Enable = true
    }
}

// Função para calcular as variáveis intermediárias
func calcularVariaveisIntermediarias(offer *models.Offer) {
    offer.R1 = offer.Discount / offer.Price
    offer.R2 = offer.Commission / offer.Price
    offer.R3 = float64(offer.UnitsSold) / float64(offer.PromoUnits)

    now := time.Now()
    durationTotal := offer.EndDate.Sub(offer.StartDate)
    durationRemaining := offer.EndDate.Sub(now)
    offer.R4 = transfCountDateTime(durationRemaining) / transfCountDateTime(durationTotal)
}

// Função para calcular o valor de y com base nas variáveis intermediárias e coeficientes
func calcularY(a, b, c, d float64, r1, r2, r3, r4 float64) float64 {
    return a*r1 + b*r2 + c*r3 + d*r4
}

// Função para normalizar os valores de y
func normalizarYs(offers []models.Offer) {
    var minY, maxY float64
    minY, maxY = offers[0].YNorm, offers[0].YNorm

    for _, offer := range offers {
        if offer.YNorm < minY {
            minY = offer.YNorm
        }
        if offer.YNorm > maxY {
            maxY = offer.YNorm
        }
    }

    for i := range offers {
        offers[i].YNorm = (offers[i].YNorm - minY) / (maxY - minY)
    }
}

// Função para ordenar as ofertas
func ordenarOfertas(offers []models.Offer, a, b, c, d float64) []models.Offer {
    // Calculando y para cada oferta
    for i := range offers {
        verificarValores(&offers[i])
        if offers[i].Enable {
            calcularVariaveisIntermediarias(&offers[i])
            offers[i].YNorm = calcularY(a, b, c, d, offers[i].R1, offers[i].R2, offers[i].R3, offers[i].R4)
        }
    }

    // Filtrando apenas as ofertas habilitadas
    ofertasHabilitadas := []models.Offer{}
    for _, offer := range offers {
        if offer.Enable {
            ofertasHabilitadas = append(ofertasHabilitadas, offer)
        }
    }

    // Normalizando os valores de y
    normalizarYs(ofertasHabilitadas)

    // Ordenando as ofertas pelas 10 melhores (maior y_norm)
    sort.SliceStable(ofertasHabilitadas, func(i, j int) bool {
        return ofertasHabilitadas[i].YNorm > ofertasHabilitadas[j].YNorm
    })

    // Mantendo as 10 primeiras fixas e embaralhando as demais
    top10 := ofertasHabilitadas[:10]
    restantes := ofertasHabilitadas[10:]

    rand.Seed(time.Now().UnixNano())
    rand.Shuffle(len(restantes), func(i, j int) {
        restantes[i], restantes[j] = restantes[j], restantes[i]
    })

    return append(top10, restantes...)
}

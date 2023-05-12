package config

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/ShikharY10/gbAUTH/cmd/models"
	"github.com/ShikharY10/gbAUTH/cmd/utils"
	"github.com/cloudinary/cloudinary-go"
	"github.com/cloudinary/cloudinary-go/api/uploader"
)

type Cloudinary struct {
	CloudName    string
	APIKey       string
	APISecret    string
	AvatarFolder string
	cloudinary   *cloudinary.Cloudinary
}

func InitCloudinary(env *ENV) *Cloudinary {

	var cloud Cloudinary
	cld, err := cloudinary.NewFromParams(
		env.CLOUDINARY_CLOUD_NAME,
		env.CLOUDINARY_API_KEY,
		env.CLOUDINARY_API_SECRET,
	)
	if err != nil {
		panic(err)
	}

	cloud.APIKey = env.CLOUDINARY_API_KEY
	cloud.APISecret = env.CLOUDINARY_API_SECRET
	cloud.CloudName = env.CLOUDINARY_CLOUD_NAME
	cloud.AvatarFolder = env.CLOUDINARY_AVATAR_FOLDER_NAME
	cloud.cloudinary = cld

	return &cloud
}

func (cloud *Cloudinary) UploadUserAvatar(tempName string, imageData string, extension string) (*models.Avatar, error) {
	var image []byte = utils.Decode(imageData)
	f, err := os.Create("temp/" + tempName + "." + extension)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	_, err = f.Write(image)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()
	uploadParam, err := cloud.cloudinary.Upload.Upload(
		ctx,
		"temp/"+tempName+"."+extension,
		uploader.UploadParams{Folder: cloud.AvatarFolder},
	)
	if err != nil {
		return nil, err
	}
	os.Remove("temp/" + tempName + "." + extension)

	var avatar models.Avatar
	avatar.PublicId = uploadParam.PublicID
	avatar.SecureUrl = uploadParam.SecureURL
	return &avatar, nil
}

func (cloud *Cloudinary) DeleteUserAvatar(publicId string) error {
	fmt.Println("publicId: ", publicId)
	param := uploader.DestroyParams{
		PublicID:   publicId,
		Invalidate: true,
	}
	result, err := cloud.cloudinary.Upload.Destroy(
		context.TODO(),
		param,
	)
	fmt.Println("response code: ", result.Response.StatusCode)
	fmt.Println("result: ", result.Result)
	if err != nil {
		return err
	} else {
		return nil
	}
}

// func ImageUploadHelper(input interface{}) (string, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()

// 	//create cloudinary instance
// 	cld, err := cloudinary.NewFromParams("", "", "")
// 	if err != nil {
// 		return "", err
// 	}

// 	//upload file
// 	uploadParam, err := cld.Upload.Upload(ctx, input, uploader.UploadParams{Folder: "config.EnvCloudUploadFolder()"})
// 	if err != nil {
// 		return "", err
// 	}
// 	return uploadParam.SecureURL, nil
// }
